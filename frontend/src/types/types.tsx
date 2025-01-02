export type Post = {
  id: number,
  title: string,
  content: string,
  created_at: number,
  updated_at: number
}

export type Error = {
  isError: boolean,
  message: string
}