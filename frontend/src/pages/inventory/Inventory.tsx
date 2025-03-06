import CategoryList from "./CategoryList";
import { createSignal, Match, Switch } from "solid-js";
import "../../styles/Inventory.css";

const tabs = [
    { name: "Data" },
    { name: "Products" },
    { name: "Categories" } 
]

export default function Inventory() {
    const [selectedTab, setSelectedTab] = createSignal(0);

    return (
        <main>
            <h1>Inventory Management</h1>
            <br />
            <nav>
                <ul class="inventory-nav">
                    {tabs.map((tab, i) => (
                        <li>
                            <button onClick={() => setSelectedTab(i)} 
                                class={selectedTab() === i ? "selected" : ""}>
                                    {tab.name}
                            </button>
                        </li>))
                    }
                </ul>
            </nav>
            <Switch>
                <Match when={selectedTab() === 0}>
                    <div>Products</div>
                </Match>
                <Match when={selectedTab() === 1}>
                    <div>Categories</div>
                </Match>
                <Match when={selectedTab() === 2}>
                    <CategoryList />
                </Match>
            </Switch>
        </main>
    )
}