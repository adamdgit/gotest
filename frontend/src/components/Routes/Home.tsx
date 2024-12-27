import React, { useState } from 'react'

export default function Home() {

  const [count, setCount] = useState(0);
  console.log("re-render")

  return (
    <div>
      <p>Count is: {count}</p>
      <button onClick={() => setCount(count +1)}>Increment</button>
    </div>
  )
}
