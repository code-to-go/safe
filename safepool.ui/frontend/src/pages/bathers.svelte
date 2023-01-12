<script>
  import {
    BlockTitle,
    Button,
    f7,
    Icon,
    Link,
    List,
    ListButton,
    ListInput,
    ListItem,
    Navbar,
    NavLeft,
    NavRight,
    NavTitle,
    Page,
    Searchbar,
    Subnavbar,
    theme,
  } from "framework7-svelte";
  import ListIndexComponent from "framework7/components/list-index";
  import { onMount } from "svelte";
  import { UpdateIdentity, GetIdentities } from "../../wailsjs/go/main/App";
  import Addbather from './addbather.svelte'

  export let poolName;
  let identities = [];

  onMount(async () => {
    identities = (await GetIdentities(poolName)) || [];
    console.log(identities);
  });

  async function updateNick(identity, e) {
    if (e.target.value) {
      const i = {...identity, n: e.target.value}
      await UpdateIdentity(i)
    }
  }

  let popupOpened = false
</script>

<Page>
  <Addbather poolName={poolName} {popupOpened} onPopupClosed={() => popupOpened = false}/>
  <List mediaList ul={false}>
    <ul>
      <ListButton onClick={()=>{popupOpened = true}}>Add</ListButton>
      {#each identities as identity, index (index)}
        <ListInput
          value={identity.n}
          onChange={e=>updateNick(identity, e)}
          info={identity.s.u}
        />
      {/each}
    </ul>
  </List>
</Page>

<style>
</style>
