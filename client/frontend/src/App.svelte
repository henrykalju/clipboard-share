<script lang="ts">
  import {EventsOn} from '../wailsjs/runtime'
  import {GetHistory, WriteToCB} from '../wailsjs/go/main/App'
  import {common} from '../wailsjs/go/models'

  let items: common.ItemWithID[] = $state([])
  GetHistory().then(result => items = result.reverse())

  function cbUpdate(...data: any): void {
    console.log("cb updated")
    GetHistory().then(result => items = result.reverse())
  }

  function handleClick(item: common.ItemWithID) {
    WriteToCB(item.ID)
  }
  
  EventsOn("CB_UPDATE_EVENT", cbUpdate)
</script>

<main>
  <ul>
    {#each items as item}
      <li key={item.ID} on:click={() => handleClick(item)}>
        {item.Text}
        <span class="tooltip-text">{item.Text}</span>
      </li>
    {/each}
  </ul>
</main>

<style>
  ul {
    list-style-type: none;
    padding: 0;
    margin: 0;
  }

  li {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 300px; /* Adjust as needed */
    padding: 5px;
    border-bottom: 1px solid #ccc;
    cursor: pointer;
    position: relative;
  }

  /* Tooltip Text */
  .tooltip-text {
    visibility: hidden;
    opacity: 0;
    background-color: rgba(0, 0, 0, 0.8);
    color: #fff;
    text-align: center;
    padding: 5px;
    border-radius: 4px;
    white-space: pre-wrap;
    position: absolute;
    bottom: 120%;
    left: 50%;
    transform: translateX(-50%);
    min-width: 300px;
    word-break: break-word;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    z-index: 100;
    transition: opacity 0.2s ease-in-out;
    pointer-events: none;
  }

  /* Show tooltip on hover */
  li:hover .tooltip-text {
    visibility: visible;
    opacity: 1;
  }
</style>