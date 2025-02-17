import { createSignal, Show } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Product } from "../../types/types";
import CategoryList from "./CategoryList";
import "../../styles/Products.css";

export default function Inventory() {
    const { userData } = useAuth();

    const [products, setProducts] = createSignal<Product[]>([]);
    const [loading, setLoading] = createSignal(false);
    const [error, setError] = createSignal(false);

    return (
        <main>
            <h1>Inventory Management</h1>
            <Show when={error()}>
                <p style={{color: 'red'}}>Something went wrong</p> 
            </Show>

            <CategoryList />
        </main>
    )
}