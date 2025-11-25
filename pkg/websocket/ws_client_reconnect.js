// WebSocket Client with Reconnection Logic (JavaScript)
// For AdminPanel / CustomShell web interface

class WebSocketClient {
    constructor(url, token) {
        this.url = url;
        this.token = token;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectDelay = 30000; // 30 seconds
        this.baseReconnectDelay = 1000; // 1 second
        this.handlers = new Map();
    }

    connect() {
        const wsUrl = `${this.url}?token=${this.token}`;
        console.log('Connecting to WebSocket:', wsUrl);

        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.reconnectAttempts = 0; // Reset on successful connection
            this.onOpen();
        };

        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                this.handleMessage(message);
            } catch (error) {
                console.error('Failed to parse message:', error);
            }
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.ws.onclose = (event) => {
            console.log('WebSocket closed:', event.code, event.reason);
            this.scheduleReconnect();
        };
    }

    scheduleReconnect() {
        // Exponential backoff: delay = min(baseDelay * 2^attempts, maxDelay)
        const delay = Math.min(
            this.baseReconnectDelay * Math.pow(2, this.reconnectAttempts),
            this.maxReconnectDelay
        );

        this.reconnectAttempts++;
        console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})...`);

        setTimeout(() => {
            this.connect();
        }, delay);
    }

    send(type, payload, agentId = null) {
        if (this.ws.readyState !== WebSocket.OPEN) {
            console.error('WebSocket is not connected');
            return;
        }

        const message = {
            type: type,
            payload: payload,
            timestamp: new Date().toISOString(),
            agent_id: agentId,
            request_id: this.generateRequestId()
        };

        this.ws.send(JSON.stringify(message));
    }

    on(messageType, handler) {
        this.handlers.set(messageType, handler);
    }

    handleMessage(message) {
        const handler = this.handlers.get(message.type);
        if (handler) {
            handler(message);
        } else {
            console.warn('No handler for message type:', message.type);
        }
    }

    onOpen() {
        // Override in subclass or set externally
    }

    generateRequestId() {
        return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }
}

// Usage Example:
// const client = new WebSocketClient('ws://localhost:8000/ws', 'your-jwt-token');
// client.on('agent_status', (message) => {
//     console.log('Agent status:', message.payload);
// });
// client.connect();