import { For } from "solid-js";
import "../styles/Skeleton.css"

export default function SkeletonLoader({ items } : { items: number }) {

  return (
    <div class="loader-wrap">
        <ul class="skeleton-list">
            <For each={Array.from({length: items})}>
                {item => (
                    <li></li>
                )}
            </For>
        </ul>
    </div>
  )
}
