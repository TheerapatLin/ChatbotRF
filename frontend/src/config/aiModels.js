// AI Provider และ Model Configuration
export const AI_PROVIDERS = {
  OPENAI: 'openai',
  BEDROCK: 'bedrock'
}

export const OPENAI_MODELS = [
  {
    id: 'gpt-4o',
    name: 'GPT-4o',
    description: 'Latest GPT-4 Omni model with enhanced multimodal capabilities',
    contextWindow: '128K tokens',
    pricing: '$5/$15 per 1M tokens',
    capabilities: ['text', 'vision'],
    recommended: true,
    category: 'gpt4'
  },
  {
    id: 'gpt-4o-mini',
    name: 'GPT-4o Mini',
    description: 'Cost-efficient small model for fast, lightweight tasks',
    contextWindow: '128K tokens',
    pricing: '$0.15/$0.60 per 1M tokens',
    capabilities: ['text', 'vision'],
    recommended: false,
    category: 'gpt4'
  },
  {
    id: 'gpt-4-turbo',
    name: 'GPT-4 Turbo',
    description: 'Previous generation high-intelligence model',
    contextWindow: '128K tokens',
    pricing: '$10/$30 per 1M tokens',
    capabilities: ['text', 'vision'],
    recommended: false,
    category: 'gpt4'
  }
]

export const BEDROCK_MODELS = [
  {
    id: 'apac.anthropic.claude-sonnet-4-5-v1:0',
    name: 'Claude Sonnet 4.5 (APAC)',
    description: 'Latest Claude Sonnet with APAC cross-region routing for better performance',
    contextWindow: '200K tokens',
    pricing: '$3/$15 per 1M tokens',
    capabilities: ['text', 'vision'],
    recommended: true,
    category: 'claude'
  },
  {
    id: 'anthropic.claude-3-5-sonnet-20241022-v2:0',
    name: 'Claude 3.5 Sonnet v2',
    description: 'Standard Claude 3.5 Sonnet (October 2024 version)',
    contextWindow: '200K tokens',
    pricing: '$3/$15 per 1M tokens',
    capabilities: ['text', 'vision'],
    recommended: false,
    category: 'claude'
  },
  {
    id: 'apac.anthropic.claude-3-5-haiku-v1:0',
    name: 'Claude Haiku (APAC)',
    description: 'Fast and efficient Claude model with APAC routing',
    contextWindow: '200K tokens',
    pricing: '$1/$5 per 1M tokens',
    capabilities: ['text', 'vision'],
    recommended: false,
    category: 'claude'
  },
  {
    id: 'apac.anthropic.claude-3-opus-v1:0',
    name: 'Claude Opus (APAC)',
    description: 'Most capable Claude model for complex tasks',
    contextWindow: '200K tokens',
    pricing: '$15/$75 per 1M tokens',
    capabilities: ['text', 'vision'],
    recommended: false,
    category: 'claude'
  },
  {
    id: 'amazon.nova-pro-v1:0',
    name: 'Amazon Nova Pro',
    description: 'Highly capable multimodal model for complex tasks',
    contextWindow: '300K tokens',
    pricing: '$0.80/$3.20 per 1M tokens',
    capabilities: ['text', 'vision', 'video'],
    recommended: false,
    category: 'nova'
  },
  {
    id: 'amazon.nova-lite-v1:0',
    name: 'Amazon Nova Lite',
    description: 'Low-cost multimodal model for processing images, videos, and text',
    contextWindow: '300K tokens',
    pricing: '$0.06/$0.24 per 1M tokens',
    capabilities: ['text', 'vision', 'video'],
    recommended: false,
    category: 'nova'
  },
  {
    id: 'amazon.nova-micro-v1:0',
    name: 'Amazon Nova Micro',
    description: 'Text-only model optimized for speed and cost',
    contextWindow: '128K tokens',
    pricing: '$0.035/$0.14 per 1M tokens',
    capabilities: ['text'],
    recommended: false,
    category: 'nova'
  },
  {
    id: 'amazon.titan-text-premier-v1:0',
    name: 'Amazon Titan Premier',
    description: 'Advanced RAG-optimized model with enterprise features',
    contextWindow: '32K tokens',
    pricing: '$0.50/$1.50 per 1M tokens',
    capabilities: ['text'],
    recommended: false,
    category: 'titan'
  }
]

export const MODEL_CATEGORIES = {
  claude: {
    name: 'Claude Models',
    description: 'Anthropic Claude models - powerful and versatile',
    icon: '🤖'
  },
  nova: {
    name: 'Amazon Nova',
    description: 'Budget-friendly multimodal models',
    icon: '⭐'
  },
  titan: {
    name: 'Amazon Titan',
    description: 'Enterprise RAG-optimized models',
    icon: '🏢'
  },
  gpt4: {
    name: 'GPT-4 Series',
    description: 'OpenAI GPT-4 family models',
    icon: '🔥'
  }
}

/**
 * Get all models for a specific provider
 * @param {string} provider - 'openai' or 'bedrock'
 * @returns {Array} Array of model objects
 */
export function getModelsByProvider(provider) {
  return provider === AI_PROVIDERS.OPENAI ? OPENAI_MODELS : BEDROCK_MODELS
}

/**
 * Get the recommended model for a provider
 * @param {string} provider - 'openai' or 'bedrock'
 * @returns {Object} Recommended model object
 */
export function getRecommendedModel(provider) {
  const models = getModelsByProvider(provider)
  return models.find(m => m.recommended) || models[0]
}

/**
 * Get models by category
 * @param {string} category - 'claude', 'nova', 'titan', or 'gpt4'
 * @returns {Array} Array of model objects in that category
 */
export function getModelsByCategory(category) {
  const allModels = [...OPENAI_MODELS, ...BEDROCK_MODELS]
  return allModels.filter(m => m.category === category)
}

/**
 * Find a model by ID
 * @param {string} modelId - Model identifier
 * @returns {Object|null} Model object or null if not found
 */
export function findModelById(modelId) {
  const allModels = [...OPENAI_MODELS, ...BEDROCK_MODELS]
  return allModels.find(m => m.id === modelId) || null
}
