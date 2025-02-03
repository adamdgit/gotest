import { useAuth } from '../AuthProvider';
import { useNavigate } from '@solidjs/router';

export default function Home() {
  const { userData } = useAuth();
  const navigate = useNavigate();

  if (!userData.email) {
    navigate("/");
  }

  return (
    <main>
      Home page
    </main>
  )
}
