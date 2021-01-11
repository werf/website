
var socket = io();

// submit text message without reload/refresh the page
$("form").submit(function (e) {
  e.preventDefault(); // prevents page reloading
  socket.emit("chat_message", $("#txt").val());
  $("#txt").val("");
  return false;
});

// append the chat text message
socket.on("chat_message", function (msg) {
  $("#messages").append($("<li>").html(msg));
  $("#messages").animate(
    {
      scrollTop: $("#messages li").last().offset().top,
    },
    "slow"
  );
});

// append text if someone is online
socket.on("is_online", function (username) {
  $("#messages").append($("<li>").html(username));
  $("#messages").animate(
    {
      scrollTop: $("#messages li").last().offset().top,
    },
    "slow"
  );
});

// ask username
socket.emit("login");

//Emit typing
$("#txt").bind("keypress", () => {
  socket.emit("typing");
});

socket.on("typing", (data) => {
  feedback.html(data.username + " is typing a message...");
});