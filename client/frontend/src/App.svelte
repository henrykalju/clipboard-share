<script lang="ts">
  import {EventsOn} from '../wailsjs/runtime'
  import {GetHistory, WriteToCB, GetConfig, UpdateConfig} from '../wailsjs/go/main/App'
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

  let showModal = $state(false);
  let username = $state("");
  let password = $state("");
  let url = $state("");
  
  function saveConfig() {
    let c: common.Config = {
      BackendUrl: url,
      Username: username,
      Password: password
    };

    UpdateConfig(c).catch(err => alert(err));
    showModal = false;
  }

  async function openModal() {
    let c: common.Config = await GetConfig();
    
    url = c.BackendUrl;
    username = c.Username;
    password = c.Password;

    showModal = true;
  }
</script>

<main>
  <button onclick={openModal}>⚙️ Settings</button>
  {#if showModal}
    <div role="button" tabindex="0" class="modal-backdrop" onclick={() => showModal = false} onkeydown={(e) => {if (e.key === "Escape") showModal = false;}}>
      <div class="modal" role="dialog" tabindex="0" onclick={(e) => e.stopPropagation()} onkeydown={() => {}}>
        <h2>Settings</h2>
        <label>
          Username:
          <input bind:value={username} />
        </label>
        <label>
          Password:
          <input bind:value={password} />
        </label>
        <label>
          URL:
          <input bind:value={url} />
        </label>
        <button onclick={saveConfig}>Save</button>
        <button onclick={() => showModal = false}>Close</button>
      </div>
    </div>
  {/if}
  <ul>
    {#each items as item}
      <li key={item.ID}>
        <button style="width: 100%;" onclick={() => handleClick(item)}>
          {item.Text}
          <span class="tooltip-text">{item.Text}</span>
        </button>
      </li>
    {/each}
  </ul>
</main>

<style>
  label {
    color: black;
  }

  .modal-backdrop {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex; justify-content: center; align-items: center;
    z-index: 99;
  }
  .modal {
    background: white;
    padding: 1rem;
    border-radius: 8px;
    width: 300px;
  }

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