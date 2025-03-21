<script lang="ts">
  import {EventsOn} from '../wailsjs/runtime'

  let items = $state([
    {id: 1, text: "item 1"},
    {id: 2, text: "item 2"},
    {id: 3, text: "item 3"}
  ])

  function cbUpdate(...data: any): void {
    console.log("cb updated")
  }
  
  EventsOn("CB_UPDATE_EVENT", cbUpdate)
</script>

<main>
  <ul>
    {#each items as {id, text}}
      <li key={id}>{text}</li>
    {/each}
  </ul>
</main>

<style>

</style>
