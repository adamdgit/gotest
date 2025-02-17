import { createSignal, For, onMount } from "solid-js";
import { Category } from "../../types/types";
import  "../../styles/category.css"

export default function CategoryList() {
    const [catname, setCatname] = createSignal("")
    const [catDesc, setCatDesc] = createSignal("")
    const [categories, setCategories] = createSignal<Category[]>([]);
    
    // get categories on mount
    onMount(async () => {
        const res = await fetch('http://127.0.0.1:8081/api/v1/categories', {
            method: 'GET',
            credentials: 'include',
            headers: {
                "content-type": "application/json"
            }
        });

        if (res.ok) {
            const data = await res.json() as Category[];
            setCategories(data);
            console.log(data)
        } else {
            console.log(res.status, res.statusText)
        }
    })

    async function addNewCategory(e: Event) {
        e.preventDefault();

        const res = await fetch(`http://127.0.0.1:8081/api/v1/categories`, {
            method: 'PUT',
            credentials: 'include',
            headers: {
                "content-type": "application/json"
            },
            body: JSON.stringify({
                "name": catname(),
                "description": catDesc()
            })
        });

        if (res.ok) {

        } else {

        }
    }

  return (
    <>
        <h2>Add new category</h2>
        <form class="db-form">
            <label for="categoryName">Name</label>
            <input type="text" name="categoryName" 
                onchange={(e) => setCatname(e.target.value)} 
            />

            <label for="categoryDesc">Description</label>
            <input type="text" name="categoryDesc" 
                onchange={(e) => setCatDesc(e.target.value)} 
            />

            <button onclick={(e) => addNewCategory(e)}>Add</button>
        </form>

        <h3>Categories</h3>
        <ul class="category-list">
            <li>
                <strong>Name</strong>
                <strong>Description</strong>
            </li>
            <For each={categories()}>
                {category => (
                    <li value={category.id}>
                        <div>{category.name}</div>
                        <div>{category.description}</div>
                    </li>
                )}
            </For>
        </ul>
    </>
  )
}
