import { useState } from "react";

type Post = {
  id: number,
  title: string,
  content: string,
  created_at: number,
  updated_at: number
}

export function useFetch(url: string) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<Post[] | []>([]);
  
  const fetchPosts = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(url);
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

  return { data, loading, error, fetchPosts };
}