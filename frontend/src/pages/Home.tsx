import { useAuth } from '../AuthProvider';

export default function Home() {
  const { userData } = useAuth();

  return (
    <>
      <header>
          <nav class="nabvar">
              <ul class="navlist">
                  <li><a href="/home">Home</a></li>
                  <li><a href="/">Login</a></li>
                  <li><a href="/register">Register</a></li>
              </ul>
          </nav>
          <div class="profilewrapper">
            {userData?.email ?? "guest"}
            {userData ? (
                <button class='btn'>Logout</button>
              ) : (
              <a href="/" class='btn'>Login</a>
            )}
          </div>
      </header>
      <main>
        Home page
      </main>
    </>
  )
}
