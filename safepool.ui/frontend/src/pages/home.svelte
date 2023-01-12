<script>
  import {
    Page,
    Navbar,
    NavLeft,
    NavTitle,
    NavTitleLarge,
    NavRight,
    Link,
    Toolbar,
    Block,
    BlockTitle,
    List,
    ListItem,
    Row,
    Col,
    Button,
    Icon,
    Views,
    View,
    Tabs,
    Tab,
    theme,
    Panel,
    Popup,
    ListInput,
  } from "framework7-svelte";
  import { onMount } from "svelte";

  import { GetPoolList, GetSelf } from "../../wailsjs/go/main/App";
    import Addpool from "./addpool.svelte";
  import Addpopup from "./addpool.svelte";

  let pools = ["something not working!"];
  let selfToken

  async function getPools() {
    console.log('fetch pool list')
    pools = await GetPoolList() || [];
  }
  async function getSelfToken() {
    selfToken = btoa(await GetSelf())
  }
  let popupOpened = false;
  const style = "height: 8em";
</script>

<Addpool {popupOpened} onPopupClosed={() => popupOpened = false}/>

<Page name="view-home" onPageBeforeIn={getPools}>
  <!-- Views/Tabs container -->
  <Views tabs class="pool-areas">
    <!-- Tabbar for switching views-tabs -->
    <Toolbar tabbar labels bottom>
      <Link
        tabLink="#view-pool"
        tabLinkActive
        iconIos="f7:house_fill"
        iconAurora="f7:house_fill"
        iconMd="material:home"
        text="Home"
      />
      <Link
        tabLink="#view-settings"
        iconIos="f7:gear"
        iconAurora="f7:gear"
        iconMd="material:settings"
        text="Settings"
      />
    </Toolbar>

    <Tabs>
      <Tab id="view-pool" tabActive>
        <Navbar>
          <NavTitle>Safepool <Icon material="water" color="blue" /></NavTitle>
          <NavRight>
            <Button tabletFullscreen onClick={() => (popupOpened = true)}
              >Add</Button
            >
          </NavRight>
        </Navbar>
        <List class="components-list">
          {#each pools as pool}
            <ListItem
              reloadDetail={theme.aurora}
              link="/pools/{pool.replace('/', '|')}"
              title={pool}
            >
              <i class="icon icon-f7" slot="media" />
            </ListItem>
          {/each}
        </List>
      </Tab>
      <Tab id="view-settings" onTabShow={getSelfToken}>
        <List>
          <ListInput 
          placeholder="Identity"
          type="textarea"
          inputStyle={style}
          info="your identity token"
          value={selfToken} 
          readonly/>
        </List>
      </Tab>
    </Tabs>

    <!-- Settings View -->
  </Views>
</Page>
