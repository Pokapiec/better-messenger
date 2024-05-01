<script>
  let usernameElem;
  let modalWrapper;
  import { onMount } from "svelte";

  function onSubmitUsername(e) {
    e.preventDefault();
    localStorage.setItem("username", usernameElem.value);
    hideUsernameModal();
  }

  function hideUsernameModal() {
    modalWrapper.classList.add("display-none");
  }

  onMount(() => {
    const username = localStorage.getItem("username");

    if (username !== undefined && username !== null && username !== "") {
      modalWrapper.classList.add("display-none");
    }
  });
</script>

<div class="username-modal-wrapper" bind:this={modalWrapper}>
  <div class="username-modal">
    <h1>Enter your chat username</h1>
    <form id="username-form" on:submit|preventDefault={onSubmitUsername}>
      <input
        type="text"
        name="username-input"
        id="username-input"
        placeholder="Type your username..."
        maxlength="30"
        minlength="3"
        required
        bind:this={usernameElem}
      />
      <input type="submit" value="submit" />
    </form>
  </div>
</div>
