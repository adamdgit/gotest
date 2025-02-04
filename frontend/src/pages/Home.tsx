import { useAuth } from '../AuthProvider';

export default function Home() {
  const { userData } = useAuth();

  return (
    <main>
      Home page
    </main>
  )
}
