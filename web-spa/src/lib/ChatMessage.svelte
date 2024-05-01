<script>
  export let remote;
  export let username;
  export let message;
</script>

<div class="chat-msg" class:remote-msg={remote} class:client-msg={!remote}>
  <div class="msg-content" class:client-msg-content={!remote}>
    {#if remote}
      <div class="msg-username">
        {username}
      </div>
    {/if}
    {message}
  </div>
</div>
