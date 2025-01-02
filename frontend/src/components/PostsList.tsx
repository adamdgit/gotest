import { useEffect, useState } from 'react'
import type { Error, Post } from '../types/types';

export default function PostsList() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error>({isError: false, message: ""});
  const [data, setData] = useState<Post[] | []>([]);

  const fetchPosts = async () => {
    setLoading(true);
    setError({isError: false, message: ""});
    setTimeout(async () => {
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
        console.log(err)
        setError({isError: true, message: "Error fetching posts"});
      } finally {
        setLoading(false);
      }
    }, 1000)
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  if (loading) return (
    <p>Fetching posts...</p>
  )

  if (error.isError) return (
    <p>{error.message}</p>
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
