<script>
  import {
    Navbar,
    Page,
    Block,
    Tabs,
    Tab,
    Link,
    Toolbar,
    NavRight,
    Button,
    Icon,
  } from "framework7-svelte";
  import Chat from "./chat.svelte";
  import Library from "./library.svelte";
  import People from "./bathers.svelte";
  import Bathers from "./bathers.svelte";

  export let poolName;
  let activeFeature = "chat";
</script>

<Page pageContent={false}>
  <Navbar title={poolName} backLink="Back">
    <NavRight>
      <Button
        active={activeFeature == "chat"}
        onClick={() => (activeFeature = "chat")}
      >
      <Icon material="chat" />
      Chat
      </Button>
      <Button
        active={activeFeature == "library"}
        onClick={() => (activeFeature = "library")}
      >
        <Icon material="folder_special" />
        Library
      </Button>
      <Button
        active={activeFeature == "people"}
        onClick={() => (activeFeature = "bathers")}
      >
        <Icon material="face" />
        Bathers
      </Button>
    </NavRight>
  </Navbar>
  {#if activeFeature == "chat"}
    <Chat {poolName} />
  {/if}
  {#if activeFeature == "library"}
    <Library {poolName} />
  {/if}
  {#if activeFeature == "bathers"}
    <Bathers {poolName} />
  {/if}
</Page>
