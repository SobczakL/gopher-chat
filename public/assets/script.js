console.log("hello")

class User {
    constructor(id, roomId) {
        this.id = id
        this.roomId = roomId
        this.element = this.createUserElement()
        this.socket = new WebSocket(`ws://localhost:8080/chat?room=${roomId}`)
        this.isInitialized = false
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
            this.getElement("connection").innerHTML += "Status: Connected\n"
        }
        this.socket.onclose = () => {
            this.getElement("connection").innerHTML += "Status: Disconnected\n"
        }
        this.socket.onmessage = (e) => {
            const output = document.createElement("p")
            output.setAttribute("class", "output")
            output.innerHTML = e.data
            this.getElement('messagesContainer').appendChild(output)
            this.scrollToBottom();
        }
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

        const prevInput = document.createElement("p")
        prevInput.setAttribute("class", "prevInput")
        prevInput.innerHTML = inputBox.value
        console.log(inputBox.value)
        this.getElement('messagesContainer').appendChild(prevInput)
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
        this.socket.close()
        this.element.remove()
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
        return Math.random().toString(36).substring(2, 8)
    }

    joinRoom(roomId) {
        if (this.currentRoom) {
            this.users.forEach(user => user.destroy());
            this.users.clear();
            this.mainContainer.innerHTML = '';
        }

        this.currentRoom = roomId;
        document.getElementById("currentRoom").textContent = `Current Room: ${roomId}`;

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
        const user = new User(id, this.currentRoom);
        this.users.set(id, user);
        this.mainContainer.appendChild(user.element);
        user.initialize();
        this.updateUserCount();
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
    const chatManager = new ChatManager(2)

    document.getElementById("addUser").style.display = "none";
    document.getElementById("removeUser").style.display = "none";
})
