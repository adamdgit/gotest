import { useEffect } from 'react'
import { useFetch } from '../api/useFetch';


export default function PostsList() {
  const { data, loading, error, fetchPosts } = useFetch('http://localhost:8081/api/v1/posts');

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
        <li key={post.id}>{post.title} | {post.content}</li>
      )}
    </ul>
  )
}
