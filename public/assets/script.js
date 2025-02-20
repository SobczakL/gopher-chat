console.log("hello")

class User {
    constructor(id, roomId) {
        this.id = id
        this.roomId = roomId
        this.element = this.createUserElement()
        this.retryAttempts = 0
        this.maxRetries = 3
        this.connectWebSocket()
    }

    connectWebSocket() {
        try {
            const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${wsProtocol}//${window.location.host}/chat?room=${this.roomId}&userId=${this.id}`;
            console.log(`Attempting to connect to: ${wsUrl}`);
            
            this.socket = new WebSocket(wsUrl);
            this.setupSocketHandlers();
        } catch (error) {
            console.error("WebSocket connection error:", error);
            this.handleConnectionError();
        }
    }

    handleConnectionError() {
        if (this.retryAttempts < this.maxRetries) {
            this.retryAttempts++
            console.log(`Retrying connection... Attempt ${this.retryAttempts}`)
            setTimeout(() => this.connectWebSocket(), 1000 * this.retryAttempts)
        } else {
            this.getElement("connection").innerHTML += `Failed to connect after ${this.maxRetries} attempts\n`
        }
    }

    initialize() {
        if (this.isInitialized) return
        console.log("init")
        this.setupSocketHandlers()
        this.setupEventListeners()
        this.isInitialized = true
    }

    createUserElement() {
        const userHTML = `
    <section class="user" id="user${this.id}">
        <div class="messageWindow">
            <h3>User ${this.id} (Room: ${this.roomId})</h3>
          <div class="messagesContainer" id="messagesContainer${this.id}"></div>
          <div class="inputContainer">
            <textarea
              name="inputBox"
              class="inputBox"
              id="inputBox${this.id}"
              placeholder="Type your message..."
            ></textarea>
            <button class="button" id="sendButton${this.id}">Send</button>
          </div>
          <pre class="connection" id="connection${this.id}"></pre>
        </div>
      </section>
`
        const template = document.createElement("template")
        template.innerHTML = userHTML.trim()
        return template.content.firstElementChild
    }

    setupSocketHandlers() {
        this.socket.onopen = () => {
            console.log(`WebSocket connected for room ${this.roomId}`);
            this.retryAttempts = 0;
            this.getElement("connection").innerHTML += `Status: Connected to Room ${this.roomId}\n`;
        };

        this.socket.onclose = (event) => {
            console.log(`WebSocket closed for room ${this.roomId}:`, event);
            this.getElement("connection").innerHTML += `Status: Disconnected from Room ${this.roomId}\n`;
            if (!event.wasClean) {
                this.handleConnectionError();
            }
        };

        this.socket.onerror = (error) => {
            console.error("WebSocket Error:", error);
            this.getElement("connection").innerHTML += `Error: Connection failed\n`;
        };

        this.socket.onmessage = (event) => {
            const message = event.data;
            const messageElement = document.createElement("p");
            messageElement.setAttribute("class", "message");
            
            // Check if it's an AI message
            if (message.startsWith("AI: ")) {
                messageElement.classList.add("ai-message");
            } else {
                messageElement.classList.add("user-message");
            }
            
            messageElement.innerHTML = message;
            this.getElement('messagesContainer').appendChild(messageElement);
            this.scrollToBottom();
        };
    }
    setupEventListeners() {
        this.getElement("sendButton").addEventListener("click", () => this.send())
        this.getElement("inputBox").addEventListener("keypress", (e) => {
            if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault()
                this.send()
            }
        })
    }
    send() {
        const inputBox = this.getElement("inputBox")
        if (inputBox.value.trim() === "") return;

        this.socket.send(inputBox.value)
        inputBox.value = "";
        this.scrollToBottom()
    }
    getElement(elementName) {
        return document.getElementById(`${elementName}${this.id}`)
    }
    scrollToBottom() {
        const container = this.getElement("messagesContainer")
        container.scrollTop = container.scrollHeight
    }
    destroy() {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.close()
        }
        if (this.element && this.element.parentNode) {
            this.element.parentNode.removeChild(this.element)
        }
    }
}

class ChatManager {
    constructor(minUsers = 2) {
        this.users = new Map()
        this.minUsers = minUsers
        this.mainContainer = document.getElementById("mainContainer")
        this.userCount = document.getElementById("userCount")
        this.roomId = this.generateRoomId()
        this.currentRoom = null
        this.setupRoomControls()
    }

    setupRoomControls() {
        const roomInput = document.getElementById("roomInput");
        const joinButton = document.getElementById("joinRoom");
        const createButton = document.getElementById("createRoom");
        const currentRoomDiv = document.getElementById("currentRoom");

        // Join existing room
        joinButton.addEventListener("click", () => {
            const roomId = roomInput.value.trim();
            if (roomId) {
                this.joinRoom(roomId);
            } else {
                alert("Please enter a room ID");
            }
        });

        // Create new room
        createButton.addEventListener("click", () => {
            const roomId = this.generateRoomId();
            roomInput.value = roomId;
            this.joinRoom(roomId);
        });
    }
    generateRoomId() {
        return 'room-' + Math.random().toString(36).substring(2, 8)
    }

    joinRoom(roomId) {
        if (this.currentRoom) {
            this.users.forEach(user => {
                user.socket.close()
                user.destroy()
            });
            this.users.clear();
            this.mainContainer.innerHTML = '';
        }

        this.currentRoom = roomId;
        document.getElementById("currentRoom").textContent = `Current Room: ${roomId}`;
        document.getElementById("roomInput").value = roomId;

        // Enable user controls after joining room
        document.getElementById("addUser").style.display = "inline";
        document.getElementById("removeUser").style.display = "inline";

        // Initialize minimum users for the room
        this.initializeUsers();
    }

    initializeUsers() {
        if (!this.currentRoom) {
            alert("Please join a room first");
            return;
        }

        for (let i = 1; i <= this.minUsers; i++) {
            this.addUser();
        }
    }
    addUser() {
        if (!this.currentRoom) {
            alert("Please join a room first");
            return;
        }

        const id = this.getNextUserId();
        try {
            const user = new User(id, this.currentRoom);
            this.users.set(id, user);
            this.mainContainer.appendChild(user.element);
            user.initialize();
            this.updateUserCount();
        } catch (error) {
            console.error("Error adding user:", error);
            alert("Failed to add user. Please try again.");
        }
    }

    removeUser() {
        if (this.users.size <= this.minUsers) {
            alert(`Cannot remove user. Min ${this.minUsers} users required.`)
            return;
        }
        const lastId = Math.max(...this.users.keys())
        const user = this.users.get(lastId)
        user.destroy()
        this.users.delete(lastId)
        this.updateUserCount()
    }

    getNextUserId() {
        return this.users.size > 0 ? Math.max(...this.users.keys()) + 1 : 1
    }

    updateUserCount() {
        this.userCount.textContent = `Users: ${this.users.size}`
    }
}

document.addEventListener("DOMContentLoaded", () => {
    try {
        const chatManager = new ChatManager(2);
        
        // Hide user controls initially
        document.getElementById("addUser").style.display = "none";
        document.getElementById("removeUser").style.display = "none";

        // Add error handlers for the room controls
        document.getElementById("joinRoom").addEventListener("click", (e) => {
            const roomId = document.getElementById("roomInput").value.trim();
            if (!roomId) {
                alert("Please enter a room ID");
                return;
            }
            try {
                chatManager.joinRoom(roomId);
            } catch (error) {
                console.error("Error joining room:", error);
                alert("Failed to join room. Please try again.");
            }
        });
    } catch (error) {
        console.error("Error initializing chat:", error);
        alert("Failed to initialize chat. Please refresh the page.");
    }
});
