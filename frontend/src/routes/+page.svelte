<script lang="ts">
  let { data } = $props();
</script>

<h1>Home</h1>
<p>Logged in, protected route</p>

<p>Hello {data.user}</p>