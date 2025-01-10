<script>
	let { children } = $props();
</script>

<nav>
  <ul>
    <li data-sveltekit-preload-data="false"><a href="/">Home</a></li>
    <li><a href="/login">Login</a></li>
    <li><a href="/register">Register</a></li>
  </ul>
</nav>

<main>
  {@render children()}
</main>