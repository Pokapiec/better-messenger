<script>
  const username = localStorage.getItem("username");
</script>

<h4 class="username-badge">User: {username}</h4>

<style>
  .username-badge {
    position: fixed;
    top: 20px;
    right: 40px;
  }
</style>
