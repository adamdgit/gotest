const posts_list = document.querySelector('.posts-list');

async function get_posts() {
  const res = await fetch('http://localhost:8080');

  if (res.ok) {
    const data = await res.json();
    console.log(data)

    data.forEach(post => {
      const li = document.createElement('li');
      li.innerHTML = `${post.Id} | ${post.title} | ${post.content}`
      posts_list.appendChild(li);
    })
  } else {
    console.error(res.status, res.statusText)
  }
}

get_posts();