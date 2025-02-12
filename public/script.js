console.log("hello")
const messagesContainer = document.getElementById("messagesContainer")
const inputBox = document.getElementById("inputBox")
const connection = document.getElementById("connection")
const socket = new WebSocket("ws://localhost:8080/chat")

socket.onopen = () => {
    connection.innerHTML += "Status: Connected\n"
}
socket.onmessage = (e) => {
    const output = document.createElement("p")
    output.setAttribute("class", "output")
    output.innerHTML = e.data

    messagesContainer.appendChild(output)
}

socket.onclose = () => {
    connection.innerHTML += "Status: Disconnected"
}

const send = () => {
    if (inputBox.value.trim() === "") return;
    const prevInput = document.createElement("p")
    prevInput.setAttribute("class", "prevInput")
    prevInput.innerHTML = inputBox.value
    messagesContainer.appendChild(prevInput)
    socket.send(inputBox.value)
    inputBox.value = "";
}
