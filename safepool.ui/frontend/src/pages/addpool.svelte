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
    } from "framework7-svelte";
    import { onMount } from "svelte";

    import {
        GetConfigTemplate,
        AddPool,
        CreatePool,
    } from "../../wailsjs/go/main/App";

    export let popupOpened;
    export let onPopupClosed;

    const tokenStyle = "height: 10em";
    const configStyle = "height: 14em";
    let config;
    let token;
    let infoConfig = "Add the exchange configuration";
    let infoToken = "Paste the token you reveived from the Pool host";

    async function suggestTemplate() {
        config = await GetConfigTemplate();
    }

    async function setToken() {
        try {
            await AddPool(token);
            infoToken = "Added";
        } catch (err) {
            infoToken = "Invalid Token: " + err;
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
            <NavTitle>Add a Pool</NavTitle>
            <NavRight>
                <Link popupClose>Close</Link>
            </NavRight>
        </Navbar>
        <Toolbar tabbar labels bottom>
            <Link
                tabLink="#claim"
                tabLinkActive
                iconIos="f7:house_fill"
                iconAurora="f7:house_fill"
                iconMd="material:token"
                text="Token"
            />
            <Link
                tabLink="#create"
                iconIos="f7:gear"
                iconAurora="f7:gear"
                iconMd="material:library_add"
                text="Create"
            />
        </Toolbar>
        <Tabs>
            <Tab id="claim" tabActive>
                <List>
                    <ListInput
                        placeholder="Token"
                        inputStyle={tokenStyle}
                        info={infoToken}
                        type="textarea"
                        onChange={e => {token = e.target.value}}
                        value={token}
                    />
                    <Button onClick={setToken} fill>Add</Button>
                </List>
            </Tab>
            <Tab id="create">
                <List>
                    <ListInput
                        inputStyle={configStyle}
                        placeholder="Config"
                        info={infoConfig}
                        type="textarea"
                        onChange={e => {config = e.target.value}}
                        value={config}
                    />
                    <Block strong>
                        <Row>
                            <Col>
                                <Button onClick={suggestTemplate} fill
                                    >Template</Button
                                >
                            </Col>
                            <Col>
                                <Button onClick={setConfig} fill>Create</Button>
                            </Col>
                        </Row>
                    </Block>
                </List>
            </Tab>
        </Tabs>
    </Page>
</Popup>
