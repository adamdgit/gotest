import { createContext, useContext } from "solid-js";
import { createStore } from "solid-js/store";

export type UserData = {
    email:       string | null,
    role:        string | null,
    profile_url: string | null
}

export interface AuthContextType {
    userData:       UserData;
    validateUser:   () => Promise<boolean>
    signIn:         (email: string, password: string) => Promise<boolean>;
    signOut:        () => Promise<boolean>;
}

export const AuthContext = createContext<AuthContextType>();

export function AuthProvider(props) {
    const [userData, setUserData] = createStore<UserData>({
        email: null, role: null, profile_url: null
    })

    async function validateUser() { 
        const res = await fetch("http://127.0.0.1:8081/api/auth/getUser", {
            method: "GET",
            credentials: "include", // Ensure cookies are sent
            headers: { 
                "content-type": "application/json" 
            }
        });

        if (res.ok) {
            const { user } = await res.json();
            setUserData({ email: user.email, role: user.role, profile_url: user.profile_url })
            return true
        } else {
            console.log("not ok")
            setUserData({email: null, role: null, profile_url: null})
            return false
        }
    };

    async function signIn(email: string, password: string) {    
        const res = await fetch("http://127.0.0.1:8081/api/auth/login", {
            method: 'POST',
            credentials: 'include',
            headers: {
                "content-type": "application/json"
            }, 
            body: JSON.stringify({
                "email": email,
                "password": password,
            })
        });

        if (res.ok) {
            const { user } = await res.json();
            console.log("res data:",user)
            setUserData({ email: user.email, role: user.role, profile_url: user.profile_url })
            return true
        } else {
            console.log("Error:", res.status, res.statusText);
            setUserData({email: null, role: null, profile_url: null})
            return false
        }
    };

    async function signOut() {
        const res = await fetch("http://127.0.0.1:8081/api/auth/logout", {
            method: 'GET',
            credentials: 'include',
            headers: {
                "content-type": "application/json",
            }, 
        });

        if (res.ok) {
            setUserData({email: null, profile_url: null});
            return true
        } 
        // handle error
        else {
            console.log("Error", res.status, res.statusText);
            return false
        }
    };

  return (
    <AuthContext.Provider value={{ userData, validateUser, signIn, signOut }}>
      {props.children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("Use auth context missing")
    }
    return context;
}