<script>
  const data = fetch("http://127.0.0.1:3001/conversations");
</script>

<main>
  <section class="conversation-list">
    <h1 class="page-title">Conversations</h1>
    <div class="conversations-container">
      {#await data}
        <div>Fetching conversations...</div>
      {:then reponse}
        {#await reponse.json() then response}
          {#each response.data || [] as convo}
            <a class="conversation" href="/#/coversations/{convo.id}">
              <p>{convo.id}</p>
              <p>{convo.name}</p>
            </a>
          {/each}
        {:catch error}
          <p>error {error}</p>
        {/await}
      {/await}
    </div>
  </section>
</main>

<style>
  .conversation-list {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    margin: auto;
  }

  .page-title {
    margin: auto;
    margin-bottom: 50px;
    width: 90%;
    text-align: center;
  }

  .conversation {
    display: flex;
    flex-direction: row;
    gap: 10px;
    border: 1px solid rgb(166, 166, 166);
    border-radius: 10px;
    padding: 20px 20px;
    transition: all 0.05s ease-in-out;
    color: rgb(220, 220, 220);
    width: 100%;
    justify-content: center;
  }

  .conversations-container {
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: 90%;
    margin: auto;
    align-items: center;
  }

  .conversation:hover {
    box-shadow: 8px 8px 24px 0px rgba(66, 68, 90, 1);
    cursor: pointer;
  }
</style>
