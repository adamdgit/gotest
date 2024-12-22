import { useEffect, useState } from 'react'

type Post = {
  Id: number;         // Corresponds to `ID` in Go struct
  title: string;      // Corresponds to `Title`
  content: string;    // Corresponds to `Content`
  createdAt: string;  // Corresponds to `Created_At`, typically a string in ISO 8601 format
};

function App() {
  const [count, setCount] = useState(0);
  const [posts, setPosts] = useState<Post[]>([]);

  async function get_posts() {
    const res = await fetch("http://localhost:8080");

    if (res.ok) {
      const data = await res.json()
      setPosts(data)
    } else {
      console.error(res.status, res.statusText)
    }
  }

  // Empty dependency array means run on mount, after dom loaded
  // This runs once only
  useEffect(() => {
    get_posts();
  },[])

  return (
    <main>
      <h1>React, TS, Golang</h1>
      Count: {count}
      <br/>
      <button className='btn' onClick={() => setCount(prev => prev +1)}>Increment ++</button>

      <ul className='post-list'>
        {posts.map(post => 
          <li key={post.Id}>{post.title} | {post.content}</li>
        )}
      </ul>
    </main>
  )
}

export default App
