package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
	"github.com/sashabaranov/go-openai"
	"github.com/xuri/excelize/v2"
)

// FileService handles file analysis operations
type FileService struct {
	openaiClient   *openai.Client
	contextService *ContextService
}

// NewFileService creates a new file service
func NewFileService(client *openai.Client, contextService *ContextService) *FileService {
	return &FileService{
		openaiClient:   client,
		contextService: contextService,
	}
}

// FileAnalysisRequest represents a file analysis request
type FileAnalysisRequest struct {
	File         *multipart.FileHeader
	AnalysisType string // summary, detail, qa, extract
	Prompt       string
	Language     string // th, en
	SystemPrompt string // Optional custom system prompt
	SessionID    string // Session ID for conversation history
	UseHistory   bool   // Include conversation history in analysis
}

// FileAnalysisResponse represents the analysis response
type FileAnalysisResponse struct {
	FileID      string    `json:"file_id"`
	FileName    string    `json:"filename"`
	FileType    string    `json:"file_type"`
	FileSize    int64     `json:"file_size"`
	Analysis    string    `json:"analysis"`
	Summary     string    `json:"summary"`
	KeyPoints   []string  `json:"key_points"`
	Entities    []string  `json:"entities,omitempty"`
	Sentiment   string    `json:"sentiment,omitempty"`
	Language    string    `json:"language"`
	TokensUsed  int       `json:"tokens_used"`
	ProcessTime float64   `json:"process_time_ms"`
	Timestamp   time.Time `json:"timestamp"`
}

// Supported file types with their MIME types and max sizes
var supportedFileTypes = map[string]struct {
	mimeTypes []string
	maxSize   int64
}{
	"text": {
		mimeTypes: []string{"text/plain", "text/markdown"},
		maxSize:   10 * 1024 * 1024, // 10 MB
	},
	"document": {
		mimeTypes: []string{"application/pdf"},
		maxSize:   25 * 1024 * 1024, // 25 MB
	},
	"office": {
		mimeTypes: []string{
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",   // .docx
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",         // .xlsx
			"application/vnd.openxmlformats-officedocument.presentationml.presentation", // .pptx
		},
		maxSize: 25 * 1024 * 1024, // 25 MB
	},
	"image": {
		mimeTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
		maxSize:   20 * 1024 * 1024, // 20 MB
	},
	"code": {
		mimeTypes: []string{
			"application/javascript",
			"text/x-python",
			"text/x-go",
			"text/x-java",
		},
		maxSize: 5 * 1024 * 1024, // 5 MB
	},
	"data": {
		mimeTypes: []string{"application/json", "application/xml", "text/csv"},
		maxSize:   10 * 1024 * 1024, // 10 MB
	},
}

// ValidateFile validates file type and size
func (s *FileService) ValidateFile(file *multipart.FileHeader) error {
	// Check if file is empty
	if file.Size == 0 {
		return fmt.Errorf("file is empty")
	}

	// Get content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		// Try to detect from extension
		contentType = detectContentType(file.Filename)
	}

	// Check if file type is supported and within size limits
	var maxSize int64 = 0
	found := false

	for _, fileTypeInfo := range supportedFileTypes {
		for _, mimeType := range fileTypeInfo.mimeTypes {
			if strings.Contains(contentType, mimeType) || contentType == mimeType {
				maxSize = fileTypeInfo.maxSize
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	// Also check by extension if MIME type not found
	if !found {
		ext := getFileExtension(file.Filename)
		if isTextFile(ext) {
			maxSize = 10 * 1024 * 1024
			found = true
		}
	}

	if !found {
		return fmt.Errorf("unsupported file type: %s (allowed: pdf, docx, xlsx, txt, png, jpg, json, etc)", contentType)
	}

	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed (%d MB)", maxSize/(1024*1024))
	}

	return nil
}

// AnalyzeFile analyzes the uploaded file
func (s *FileService) AnalyzeFile(ctx context.Context, req FileAnalysisRequest) (*FileAnalysisResponse, error) {
	startTime := time.Now()

	// Debug: Log request parameters
	fmt.Printf("\n🔍 === FILE ANALYSIS REQUEST ===\n")
	fmt.Printf("   Filename: %s\n", req.File.Filename)
	fmt.Printf("   Analysis Type: %s\n", req.AnalysisType)
	fmt.Printf("   Language: %s\n", req.Language)
	fmt.Printf("   Session ID: %s\n", req.SessionID)
	fmt.Printf("   Use History: %v\n", req.UseHistory)
	fmt.Printf("   Has System Prompt: %v\n", req.SystemPrompt != "")
	fmt.Printf("   Has Custom Prompt: %v\n", req.Prompt != "")
	fmt.Printf("================================\n\n")

	// Validate file
	if err := s.ValidateFile(req.File); err != nil {
		return nil, err
	}

	// Open file
	fileData, err := req.File.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileData.Close()

	// Extract text from file based on type
	text, err := s.extractText(fileData, req.File)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from file: %w", err)
	}

	// Debug: Log extracted text length
	fmt.Printf("📄 File: %s (Size: %d bytes)\n", req.File.Filename, req.File.Size)
	fmt.Printf("📝 Extracted text length: %d characters\n", len(text))
	if len(text) > 0 {
		// Show first 100 chars of extracted text
		preview := text
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		fmt.Printf("📋 Text preview: %s\n", preview)
	} else {
		fmt.Printf("⚠️  WARNING: Extracted text is EMPTY!\n")
	}

	// If text is too long, chunk it
	if len(text) > 15000 {
		text = text[:15000] + "\n\n... (truncated)"
		fmt.Printf("✂️  Text truncated to 15000 characters\n")
	}

	// Build analysis prompt
	prompt := s.buildAnalysisPrompt(req.AnalysisType, req.Prompt, req.Language, text)
	fmt.Printf("📤 Final prompt length: %d characters\n", len(prompt))

	// Build messages array with optional system prompt
	var messages []openai.ChatCompletionMessage

	// If UseHistory is enabled and SessionID is provided, build context with conversation history
	if req.UseHistory && req.SessionID != "" && s.contextService != nil {
		// Build system prompt (use custom or default)
		systemPrompt := req.SystemPrompt
		if systemPrompt == "" {
			systemPrompt = "You are a document analysis expert. Provide clear, accurate, and structured analysis."
		}

		// Use context service to build messages with history
		var err error
		messages, err = s.contextService.BuildContextWithHistory(
			req.SessionID,
			systemPrompt,
			prompt,
			10, // last 10 messages
		)
		if err != nil {
			fmt.Printf("⚠️  Failed to build context with history: %v, falling back to simple context\n", err)
			// Fallback to simple context
			messages = s.buildSimpleContext(req.SystemPrompt, prompt)
		}
	} else {
		// Build simple context without history
		messages = s.buildSimpleContext(req.SystemPrompt, prompt)
	}

	// Call OpenAI API
	chatResp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT4oMini,
			Messages:    messages,
			Temperature: 0.7,
			MaxTokens:   2000,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to analyze with OpenAI: %w", err)
	}

	// Parse response
	analysis := chatResp.Choices[0].Message.Content

	// Extract key points (simple parsing)
	keyPoints := extractKeyPoints(analysis)

	// Calculate process time
	processTime := float64(time.Since(startTime).Milliseconds())

	// Build response
	response := &FileAnalysisResponse{
		FileID:      generateFileID(),
		FileName:    req.File.Filename,
		FileType:    req.File.Header.Get("Content-Type"),
		FileSize:    req.File.Size,
		Analysis:    analysis,
		Summary:     extractSummary(analysis),
		KeyPoints:   keyPoints,
		Language:    req.Language,
		TokensUsed:  chatResp.Usage.TotalTokens,
		ProcessTime: processTime,
		Timestamp:   time.Now(),
	}

	return response, nil
}

// extractText extracts text from different file types
func (s *FileService) extractText(file io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	contentType := fileHeader.Header.Get("Content-Type")
	ext := getFileExtension(fileHeader.Filename)

	// Read file data into memory for multiple parsing attempts
	fileData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file data: %w", err)
	}

	// Handle different file types
	switch {
	case strings.Contains(contentType, "application/pdf") || ext == ".pdf":
		return s.extractPDFText(bytes.NewReader(fileData), int64(len(fileData)))

	case strings.Contains(contentType, "application/vnd.openxmlformats-officedocument.wordprocessingml.document") || ext == ".docx":
		return s.extractDocxText(bytes.NewReader(fileData), int64(len(fileData)))

	case strings.Contains(contentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") || ext == ".xlsx":
		return s.extractExcelText(bytes.NewReader(fileData))

	case strings.Contains(contentType, "application/vnd.openxmlformats-officedocument.presentationml.presentation") || ext == ".pptx":
		return s.extractPptxText(bytes.NewReader(fileData))

	case strings.Contains(contentType, "application/xml") || ext == ".xml":
		return s.extractXMLText(fileData)

	case strings.Contains(contentType, "text/csv") || ext == ".csv":
		return string(fileData), nil

	case strings.Contains(contentType, "application/json") || ext == ".json":
		return string(fileData), nil

	case strings.Contains(contentType, "image/"):
		// Images should be handled separately by Vision API
		return fmt.Sprintf("[Image file: %s - use Vision API for analysis]", fileHeader.Filename), nil

	case strings.Contains(contentType, "text/") || isTextFile(ext):
		// Plain text files
		return string(fileData), nil

	default:
		// Try to read as text
		return string(fileData), nil
	}
}

// AnalyzeImage analyzes an image file using OpenAI Vision API
func (s *FileService) AnalyzeImage(ctx context.Context, file io.Reader, filename string, prompt string, systemPrompt string) (string, error) {
	// Read image data
	imageData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	// Encode to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// Detect image format
	imageFormat := "jpeg"
	if strings.HasSuffix(strings.ToLower(filename), ".png") {
		imageFormat = "png"
	} else if strings.HasSuffix(strings.ToLower(filename), ".gif") {
		imageFormat = "gif"
	} else if strings.HasSuffix(strings.ToLower(filename), ".webp") {
		imageFormat = "webp"
	}

	// Build prompt
	if prompt == "" {
		prompt = "Describe this image in detail. What do you see?"
	}

	// Build messages array with optional system prompt
	messages := []openai.ChatCompletionMessage{}

	// Add system prompt if provided
	if systemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	// Add user message with image
	messages = append(messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser,
		MultiContent: []openai.ChatMessagePart{
			{
				Type: openai.ChatMessagePartTypeText,
				Text: prompt,
			},
			{
				Type: openai.ChatMessagePartTypeImageURL,
				ImageURL: &openai.ChatMessageImageURL{
					URL: fmt.Sprintf("data:image/%s;base64,%s", imageFormat, base64Image),
				},
			},
		},
	})

	// Call Vision API
	resp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:     openai.GPT4oMini,
			Messages:  messages,
			MaxTokens: 1000,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to analyze image: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}

// buildAnalysisPrompt builds the prompt based on analysis type
func (s *FileService) buildAnalysisPrompt(analysisType, customPrompt, language, text string) string {
	var prompt strings.Builder

	// Add custom prompt if provided
	if customPrompt != "" {
		prompt.WriteString(customPrompt)
		prompt.WriteString("\n\n")
	}

	// Add analysis type instructions
	switch analysisType {
	case "summary":
		if language == "th" {
			prompt.WriteString("สรุปเนื้อหาของเอกสารนี้อย่างกระชับและชัดเจน")
		} else {
			prompt.WriteString("Provide a concise summary of this document.")
		}

	case "detail":
		if language == "th" {
			prompt.WriteString("วิเคราะห์เอกสารนี้อย่างละเอียด รวมถึง:\n- สาระสำคัญ\n- ประเด็นหลัก\n- ข้อมูลสำคัญ\n- ข้อสรุป")
		} else {
			prompt.WriteString("Provide a detailed analysis including:\n- Main points\n- Key information\n- Important details\n- Conclusions")
		}

	case "qa":
		if language == "th" {
			prompt.WriteString("สร้างคำถามและคำตอบที่สำคัญจากเอกสารนี้")
		} else {
			prompt.WriteString("Generate important questions and answers from this document.")
		}

	case "extract":
		if language == "th" {
			prompt.WriteString("แยกข้อมูลสำคัญจากเอกสาร เช่น ชื่อ, วันที่, ตัวเลข, สถานที่")
		} else {
			prompt.WriteString("Extract key information such as names, dates, numbers, locations from this document.")
		}

	default:
		if language == "th" {
			prompt.WriteString("วิเคราะห์และสรุปเอกสารนี้")
		} else {
			prompt.WriteString("Analyze and summarize this document.")
		}
	}

	prompt.WriteString("\n\nDocument:\n")
	prompt.WriteString(text)

	return prompt.String()
}

// buildSimpleContext builds a simple message context without history
func (s *FileService) buildSimpleContext(systemPrompt, userPrompt string) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	// Add system prompt (use custom or default)
	if systemPrompt == "" {
		systemPrompt = "You are a document analysis expert. Provide clear, accurate, and structured analysis."
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	})

	// Add user message with analysis prompt
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userPrompt,
	})

	return messages
}

// Helper functions

func generateFileID() string {
	return fmt.Sprintf("file_%d", time.Now().UnixNano())
}

func detectContentType(filename string) string {
	ext := getFileExtension(filename)
	contentTypes := map[string]string{
		".txt":  "text/plain",
		".md":   "text/markdown",
		".json": "application/json",
		".xml":  "application/xml",
		".csv":  "text/csv",
		".pdf":  "application/pdf",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".js":   "application/javascript",
		".py":   "text/x-python",
		".go":   "text/x-go",
		".java": "text/x-java",
	}

	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return "." + strings.ToLower(parts[len(parts)-1])
	}
	return ""
}

func isTextFile(ext string) bool {
	textExtensions := []string{".txt", ".md", ".json", ".xml", ".csv", ".js", ".py", ".go", ".java", ".html", ".css", ".yaml", ".yml"}
	for _, te := range textExtensions {
		if ext == te {
			return true
		}
	}
	return false
}

func extractKeyPoints(text string) []string {
	keyPoints := []string{}

	// Simple extraction: look for bullet points or numbered lists
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "•") || strings.HasPrefix(line, "*") {
			point := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "-"), "•"), "*"))
			if point != "" && len(point) > 10 {
				keyPoints = append(keyPoints, point)
			}
		}
	}

	// Limit to 5 key points
	if len(keyPoints) > 5 {
		keyPoints = keyPoints[:5]
	}

	return keyPoints
}

func extractSummary(text string) string {
	// Extract first paragraph or first 200 characters as summary
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 50 {
			if len(line) > 200 {
				return line[:200] + "..."
			}
			return line
		}
	}

	// Fallback: just take first 200 chars
	if len(text) > 200 {
		return text[:200] + "..."
	}
	return text
}

// extractPDFText extracts text from PDF files
func (s *FileService) extractPDFText(reader io.ReaderAt, size int64) (string, error) {
	pdfReader, err := pdf.NewReader(reader, size)
	if err != nil {
		fmt.Printf("❌ Failed to open PDF: %v\n", err)
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}

	var text strings.Builder
	numPages := pdfReader.NumPage()

	fmt.Printf("📖 PDF has %d pages\n", numPages)

	// Limit to first 50 pages to avoid excessive processing
	maxPages := numPages
	if maxPages > 50 {
		maxPages = 50
	}

	successfulPages := 0
	failedPages := 0

	for pageNum := 1; pageNum <= maxPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() {
			fmt.Printf("⚠️  Page %d: V is null, skipping\n", pageNum)
			failedPages++
			continue
		}

		pageText, err := page.GetPlainText(nil)
		if err != nil {
			fmt.Printf("⚠️  Page %d: GetPlainText error: %v\n", pageNum, err)
			failedPages++
			continue
		}

		if len(pageText) > 0 {
			text.WriteString(pageText)
			text.WriteString("\n\n")
			successfulPages++
			fmt.Printf("✅ Page %d: Extracted %d characters\n", pageNum, len(pageText))
		} else {
			fmt.Printf("⚠️  Page %d: No text content\n", pageNum)
			failedPages++
		}
	}

	fmt.Printf("📊 PDF extraction summary: %d successful, %d failed out of %d pages\n", successfulPages, failedPages, maxPages)

	if text.Len() == 0 {
		fmt.Printf("❌ WARNING: No text extracted from PDF! Returning placeholder.\n")
		return "[PDF file - no text content extracted]", nil
	}

	return text.String(), nil
}

// extractDocxText extracts text from Word (.docx) files
func (s *FileService) extractDocxText(reader io.ReaderAt, size int64) (string, error) {
	// Read data into memory and save to temp file (docx library needs file path)
	data := make([]byte, size)
	_, err := reader.ReadAt(data, 0)
	if err != nil {
		return "", fmt.Errorf("failed to read DOCX data: %w", err)
	}

	// Create a temporary file
	tmpFile := fmt.Sprintf("temp_docx_%d.docx", time.Now().UnixNano())
	defer func() {
		// Clean up temp file (best effort)
		_ = os.Remove(tmpFile)
	}()

	// Write to temp file
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	// Read docx file
	doc, err := docx.ReadDocxFile(tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to open DOCX: %w", err)
	}
	defer doc.Close()

	var text strings.Builder

	// Extract text content
	editable := doc.Editable()
	text.WriteString(editable.GetContent())

	if text.Len() == 0 {
		return "[Word document - no text content extracted]", nil
	}

	return text.String(), nil
}

// extractExcelText extracts text from Excel (.xlsx) files
func (s *FileService) extractExcelText(reader io.Reader) (string, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return "", fmt.Errorf("failed to open Excel: %w", err)
	}
	defer f.Close()

	var text strings.Builder

	// Get all sheet names
	sheets := f.GetSheetList()
	for _, sheetName := range sheets {
		text.WriteString(fmt.Sprintf("=== Sheet: %s ===\n", sheetName))

		// Get all rows
		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}

		for _, row := range rows {
			if len(row) > 0 {
				text.WriteString(strings.Join(row, " | "))
				text.WriteString("\n")
			}
		}

		text.WriteString("\n")
	}

	if text.Len() == 0 {
		return "[Excel file - no data extracted]", nil
	}

	return text.String(), nil
}

// extractPptxText extracts text from PowerPoint (.pptx) files
func (s *FileService) extractPptxText(reader io.Reader) (string, error) {
	// PowerPoint files are zip archives containing XML files
	// We'll read the archive and extract text from slide XML files
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read PPTX data: %w", err)
	}

	// For now, return a placeholder
	// Full PPTX parsing requires unzipping and parsing multiple XML files
	return fmt.Sprintf("[PowerPoint file - %d bytes - advanced parsing not yet implemented]", len(data)), nil
}

// extractXMLText extracts and formats XML content
func (s *FileService) extractXMLText(data []byte) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(data); err != nil {
		// If XML parsing fails, return as plain text
		return string(data), nil
	}

	var text strings.Builder

	// Pretty print XML
	doc.Indent(2)
	xmlString, err := doc.WriteToString()
	if err != nil {
		return string(data), nil
	}

	text.WriteString("=== XML Content ===\n")
	text.WriteString(xmlString)

	return text.String(), nil
}
