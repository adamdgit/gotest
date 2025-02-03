import { createSignal, onMount } from "solid-js";
import { useAuth } from "../AuthProvider";
import { For } from "solid-js";
import "../styles/Products.css";

type Product = {
    id: number,
    name: string,
    brand: string, 
    description: string,
    price: number,
}

export default function Inventory() {
    const { userData } = useAuth();

    const [products, setProducts] = createSignal<Product[] | []>([]);
    const [loading, setLoading] = createSignal(false);
    const [error, setError] = createSignal(false);

    onMount(async () => {
        setLoading(true);
        setError(false);

        const res = await fetch('http://127.0.0.1:8081/api/v1/products', {
            method: 'GET',
            credentials: 'include',
            headers: {
                "content-type": "application/json"
            },
        });

        if (res.ok) {
            const data = await res.json();
            setProducts(data);
            setLoading(false);
        }
        else {
            setError(true);
            setProducts([]);
            setLoading(false);
        }
    })

    return (
        <main>
            <h1>Product Inventory</h1>
            {error() ? <p style={{color: 'red'}}>Something went wrong</p> : null}
            {loading() ? (
                <div>Loading data...</div>
            ) : (
                <ul class="products-list">
                    <For each={products()}>
                        {product => (
                            <li>
                                <a href={`/product/${product.id}`}>
                                    <span class="prod-name">{product.name}</span>
                                    <span class="prod-desc">{product.description}</span>
                                    <span class="prod-brand">{product.brand}</span>
                                    <span class="prod-price">${product.price}</span>
                                </a>
                            </li>
                        )}
                    </For>
                </ul>
            )}
        </main>
    )
}