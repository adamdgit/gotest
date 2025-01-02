import { useEffect, useState } from 'react'

type Post = {
  id: number,
  title: string,
  content: string,
  created_at: number,
  updated_at: number
}

export default function PostsList() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<Post[] | []>([]);

  const fetchPosts = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch('http://localhost:8081/api/v1/posts', {
        credentials: "include"
      });
      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }
      const result = await res.json();
      setData(result);
    } catch (err) {
      setError((err as Error).message || 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  if (loading) return (
    <p>Fetching posts...</p>
  )

  if (error) return (
    <p>Something went wrong</p>
  )

  return (
    <ul>
      {data?.map(post => 
        <li key={post.id}>
          <h2>{post.title}</h2>
          <p>{post.content}</p>
          <span>{
            new Date(post.created_at)
              .toLocaleDateString("en-au", { minute: '2-digit', hour: '2-digit' })
          }
          </span>
        </li>
      )}
    </ul>
  )
}
