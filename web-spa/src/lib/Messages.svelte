<script>
  import { onMount } from "svelte";
  import ChatMessage from "./ChatMessage.svelte";

  export let params;

  let conn;
  let username = localStorage.getItem("username");
  let chatMsgInput;
  let messagesSpot;
  let dbMessages;

  const messagesPromise = fetch(
    `http://127.0.0.1:3001/conversations/${params.convId}/messages`
  );

  onMount(async () => {
    // const response = await fetch(
    //   `http://127.0.0.1:3001/conversations/${params.convId}/messages`
    // );
    // dbMessages = await response.json();
    // console.log(dbMessages);
    connectToWSS();
  });

  function addChatMsg(message, remote) {
    const newMessage = new ChatMessage({
      target: messagesSpot,
      props: {
        username: message.username,
        message: message.message,
        remote: remote,
      },
    });
    // Scroll to bottom in messages container
    messagesSpot.scrollTop = messagesSpot.scrollHeight;
  }

  function onMessageSubmit(e) {
    console.log("submiting form with", chatMsgInput.value);
    if (!conn) {
      return false;
    }

    if (!chatMsgInput.value) {
      return false;
    }

    let messageData = {
      username: username,
      message: chatMsgInput.value,
      conversation_id: parseInt(params.convId),
    };
    addChatMsg(messageData, false);

    let data = JSON.stringify(messageData);
    conn.send(data);
    chatMsgInput.value = "";
    return false;
  }

  function connectToWSS() {
    if (window["WebSocket"]) {
      conn = new WebSocket("ws://127.0.0.1:3001/ws");

      conn.onopen = function (evt) {
        conn.send(
          JSON.stringify({
            conversation_id: parseInt(params.convId),
            message: "<INITIAL>",
          })
        );
      };

      conn.onclose = function (evt) {
        // var item = document.createElement("div");
        // item.innerHTML = "<b>Connection closed.</b>";
        // appendLog(item);
        console.log("Closed...");
      };
      conn.onmessage = function (evt) {
        addChatMsg(JSON.parse(evt.data), true);
        console.log("Got message: ", evt.data);
      };
    } else {
      // var item = document.createElement("div");
      // item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
      // appendLog(item);
      console.log("websockets not suported...");
    }
  }
</script>

<main>
  <div>Conversation with id {params.convId}</div>
  <h1>Welcome to chatroom!</h1>
  <section id="messages" bind:this={messagesSpot}>
    {#await messagesPromise}
      <div>Fetching messages...</div>
    {:then reponse}
      {#await reponse.json() then response}
        {#each response.data || [] as msg}
          <ChatMessage
            username={msg.username}
            message={msg.message}
            remote={username !== msg.username}
          />
        {/each}
      {:catch error}
        <p>error {error}</p>
      {/await}
    {/await}

    <!-- {#each dbMessages.data as msg}
      <ChatMessage
        username={msg.username}
        message={msg.message}
        remote={username == msg.username}
      />
    {/each} -->
    <slot />
  </section>

  <form id="chat-form" on:submit|preventDefault={onMessageSubmit}>
    <input
      type="text"
      name="msg-input"
      id="msg-input"
      placeholder="Type a message..."
      bind:this={chatMsgInput}
    />
    <button type="submit">Submit</button>
  </form>
</main>
