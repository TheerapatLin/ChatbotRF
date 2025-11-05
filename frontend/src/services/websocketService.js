/**
 * WebSocket Service
 * Manages WebSocket connection for real-time chat streaming
 */
class WebSocketService {
  constructor() {
    this.ws = null
    this.messageHandlers = []
    this.errorHandlers = []
    this.closeHandlers = []
    this.openHandlers = []
    this.isConnecting = false
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectDelay = 2000 // 2 seconds
  }

  /**
   * Connect to WebSocket server
   */
  connect(url = 'ws://localhost:3001/api/chat/stream') {
    return new Promise((resolve, reject) => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        resolve()
        return
      }

      if (this.isConnecting) {
        reject(new Error('Connection already in progress'))
        return
      }

      this.isConnecting = true

      try {
        this.ws = new WebSocket(url)

        this.ws.onopen = () => {
          console.log('✅ WebSocket connected')
          this.isConnecting = false
          this.reconnectAttempts = 0
          this.openHandlers.forEach((handler) => handler())
          resolve()
        }

        this.ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data)
            this.messageHandlers.forEach((handler) => handler(data))
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error)
          }
        }

        this.ws.onerror = (error) => {
          console.error('❌ WebSocket error:', error)
          this.isConnecting = false
          this.errorHandlers.forEach((handler) => handler(error))
        }

        this.ws.onclose = (event) => {
          console.log('🔌 WebSocket disconnected')
          this.isConnecting = false
          this.closeHandlers.forEach((handler) => handler(event))

          // Auto-reconnect if not closed intentionally
          if (!event.wasClean && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++
            console.log(
              `🔄 Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`
            )
            setTimeout(() => {
              this.connect(url)
            }, this.reconnectDelay)
          }
        }
      } catch (error) {
        console.error('Failed to create WebSocket connection:', error)
        this.isConnecting = false
        reject(error)
      }
    })
  }

  /**
   * Send message through WebSocket
   * @param {Object} message - Message object to send
   */
  send(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    } else {
      console.error('WebSocket is not connected')
      throw new Error('WebSocket is not connected')
    }
  }

  /**
   * Add message handler
   * @param {Function} handler - Handler function
   */
  onMessage(handler) {
    this.messageHandlers.push(handler)
  }

  /**
   * Remove message handler
   * @param {Function} handler - Handler function to remove
   */
  offMessage(handler) {
    this.messageHandlers = this.messageHandlers.filter((h) => h !== handler)
  }

  /**
   * Add error handler
   * @param {Function} handler - Handler function
   */
  onError(handler) {
    this.errorHandlers.push(handler)
  }

  /**
   * Remove error handler
   * @param {Function} handler - Handler function to remove
   */
  offError(handler) {
    this.errorHandlers = this.errorHandlers.filter((h) => h !== handler)
  }

  /**
   * Add close handler
   * @param {Function} handler - Handler function
   */
  onClose(handler) {
    this.closeHandlers.push(handler)
  }

  /**
   * Remove close handler
   * @param {Function} handler - Handler function to remove
   */
  offClose(handler) {
    this.closeHandlers = this.closeHandlers.filter((h) => h !== handler)
  }

  /**
   * Add open handler
   * @param {Function} handler - Handler function
   */
  onOpen(handler) {
    this.openHandlers.push(handler)
  }

  /**
   * Remove open handler
   * @param {Function} handler - Handler function to remove
   */
  offOpen(handler) {
    this.openHandlers = this.openHandlers.filter((h) => h !== handler)
  }

  /**
   * Disconnect WebSocket
   */
  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
      this.reconnectAttempts = this.maxReconnectAttempts // Prevent auto-reconnect
    }
  }

  /**
   * Check if WebSocket is connected
   */
  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN
  }

  /**
   * Get connection state
   */
  getState() {
    if (!this.ws) return 'CLOSED'

    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return 'CONNECTING'
      case WebSocket.OPEN:
        return 'OPEN'
      case WebSocket.CLOSING:
        return 'CLOSING'
      case WebSocket.CLOSED:
        return 'CLOSED'
      default:
        return 'UNKNOWN'
    }
  }
}

// Export singleton instance
export default new WebSocketService()
