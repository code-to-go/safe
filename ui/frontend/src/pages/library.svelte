<script>
  import {
    AccordionContent,
    Block,
    BlockTitle,
    Button,
    Col,
    f7,
    Icon,
    Link,
    List,
    ListInput,
    ListItem,
    MenuDropdown,
    MenuDropdownItem,
    MenuItem,
    Navbar,
    NavLeft,
    NavRight,
    NavTitle,
    Page,
    Row,
    Searchbar,
    Segmented,
    Subnavbar,
    theme,
    Toggle,
  } from "framework7-svelte";
  import { onMount } from "svelte";
  import { GetPoolList, ListLibrary } from "../../wailsjs/go/main/App";
  import Chat from "./chat.svelte";

  export let poolName;
  let data = {
    items: [],
  };

  onMount(async () => {
    const documents = (await ListLibrary(poolName)) || [];
    console.log(documents);
    data.items = documents.map((d) => ({
      title: `${d.name}`,
      subtitle: `${d.localPath}`,
    }));
    console.log(data.items);
  });

  function searchAll(query, items) {
    const found = [];
    for (let i = 0; i < items.length; i += 1) {
      if (
        items[i].title.toLowerCase().indexOf(query.toLowerCase()) >= 0 ||
        query.trim() === ""
      )
        found.push(i);
    }
    return found; // return array with mathced indexes
  }

  function renderExternal(virtualList, virtualListData) {
    data = virtualListData;
  }
</script>

<Page>
  <Subnavbar inner={false}>
        <Searchbar
        searchContainer=".virtual-list"
        searchItem="li"
        searchIn=".item-title"
        clearButton
      />
  </Subnavbar>

  <List
    class="searchbar-found"
    ul={false}
    medialList
    virtualList
    virtualListParams={{
      items: data.items,
      searchAll: searchAll,
      renderExternal: renderExternal,
      height: theme.ios ? 63 : theme.md ? 73 : 77,
    }}
  >
    <ul>

      {#each data.items as item, index (index)}
        <ListItem
          mediaItem
          style={`top: ${data.topPosition}px`}
          virtualListIndex={data.items.indexOf(item)}
        >
        <span slot="title">
          <b>{item.title}</b>
        </span>
              <Row slot="after">
                <Col><Button>Open</Button></Col>
                <Col><Button>Update</Button></Col>
              </Row>
        </ListItem>
      {/each}
    </ul>
  </List>
</Page>

<style>
</style>
