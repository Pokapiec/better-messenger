<script>
  import Conversations from "./lib/Conversations.svelte";
  import Messages from "./lib/Messages.svelte";
  import Router, { location, link } from "svelte-spa-router";
  import UsernameModal from "./lib/UsernameModal.svelte";
  import UsernameBadge from "./lib/UsernameBadge.svelte";
</script>

<Router
  routes={{
    "/": Conversations,
    "/coversations/:convId": Messages,
  }}
/>
<UsernameModal />
<UsernameBadge />
