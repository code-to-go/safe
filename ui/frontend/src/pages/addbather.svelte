<script>
    import { height } from "dom7";
    import {
        Popup,
        Page,
        Toolbar,
        Link,
        Navbar,
        NavRight,
        NavTitle,
        Tab,
        Input,
        ListInput,
        List,
        Tabs,
        Block,
        Row,
        Col,
        Button,
        ListItem,
        Toggle,
    } from "framework7-svelte";
    import { onMount } from "svelte";

    import {
        GetConfigTemplate,
        AddPool,
        CreatePool,
        GetToken,
    } from "../../wailsjs/go/main/App";

    export let popupOpened;
    export let onPopupClosed;
    export let poolName

    const style = "height: 10em";
    let identity = "";
    let token = "";
    let custom = false;
    let infoConfig = "Add the exchange configuration";
    let infoIdentity = "Paste the guest's identity";

    async function suggestTemplate() {
        config = await GetConfigTemplate();
    }

    async function getToken() {
        try {
            token = await GetToken(poolName, identity);
            console.log('token:', token)
        } catch (err) {
            console.log(err)
        }
    }

    async function setConfig() {
        try {
            config = await CreatePool(config);
            infoConfig = "Created";
        } catch (err) {
            infoConfig = "Invalid Configuration: " + err;
        }
    }
</script>

<Popup tabletFullscreen opened={popupOpened} {onPopupClosed}>
    <Page>
        <Navbar>
            <NavTitle>Add Bather</NavTitle>
            <NavRight>
                <Link popupClose>Close</Link>
            </NavRight>
        </Navbar>
        <List>
            {#if token == ""}
            <ListItem>
                <span>Custom Token</span>
                <Toggle checked={custom} onChange={()=>{custom = !custom}}/>
              </ListItem>
            {#if custom}
            <ListInput
                placeholder="Public Identity"
                inputStyle={style}
                info={infoIdentity}
                type="textarea"
                onChange={(e) => {
                    identity = e.target.value;
                }}
                value={identity}
            />
            <Button onClick={getToken} fill disabled={identity.length<32}>Add Bather and get Token</Button>
            {:else}
            <Button onClick={getToken} fill>Universal Token</Button>
            {/if}
            {:else}
            <ListInput
                inputStyle={style}
                info="Generated Token"
                type="textarea"
                value={token}
                readonly
            />
            {/if}

        </List>
    </Page>
</Popup>
