<script>
  import logo from "./assets/svelte.png";
  import loader from "./assets/loader.gif";
  import Continue from "./lib/Continue.svelte";
  import DepositWidget from "./lib/DepositWidget.svelte";
  import SelectWidget from "./lib/SelectWidget.svelte";
  import axios from "axios";
  // const urlParams = new URLSearchParams(window.location.search);
  // let plugin_id = urlParams.get("plugin_id");
  let path = window.location.pathname.split("/");
  const apiUrl = import.meta.env.VITE_API_URL;
  const apiSecret = import.meta.env.VITE_API_SECRET;
  console.log(apiUrl);
  let pluginid = ""
  async function loadPlugin(path) {
    if (path[1] == "payment-gateway") {
      if (path[2] != null || path[2] != "") {
        localStorage.setItem("plugin_id", path[2])
        pluginid = path[2]
        return await axios.get(
          apiUrl+"/api/v1/public/merchant/" + path[2]
        );
      } else {
        throw Error;
      }
    }
    throw Error;
  }
  let data1 = loadPlugin(path);

  let coins = [
    { id: "0", symbol: "BTC", text: "Bitcoin", network:"bitcoin" },
    { id: "1", symbol: "USDT", text: "USD Tether", network:"ethereum" },
    { id: "2", symbol: "ETH", text: "Ethereum", network:"ethereum" },
    { id: "3", symbol: "USDC", text: "USD COin", network:"ethereum" },
  ];

  let selectedCoin;
  let amountUSD = 0.0;
  let cryptoAmount = "0.00";
  let data = {
    coin: coins[0],
    amount: {},
  };
  let pages = [SelectWidget, DepositWidget];
  let selectedPage = 0;
</script>

{#await data1}
  <div>
    <img src={loader} alt="Loader" style="width: 100%; height:100%;" />
  </div>
{:then dataA}
{#if dataA.data !== null}
  <div class="img_container">
    <!-- <img src={logo} alt="Svelte Logo" class="logo" /> -->
    <h1>{dataA.data.name}</h1>
  </div>
  <main class="main">
      {#if selectedPage == 0}
        <SelectWidget
          {coins}
          {selectedCoin}
          {data}
          {amountUSD}
          {cryptoAmount}
        />
      {/if}
      {#if selectedPage == 1}
        <DepositWidget {data} plugin_id={pluginid}/>
      {/if}
      <Continue
        bind:selectedPage
        pages={pages.length}
        {data}
        text={selectedPage == 0 ? "Continue" : "Start Over"}
      />
  </main>
  {:else}
      <div>
        <h2>An Error Occured</h2>
      </div>
    {/if}
{:catch}
  <div>
    <h2>An Error Occured</h2>
  </div>
{/await}

<style>
  :root {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
      Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
  }

  .img_container {
    display: flex;
    position: sticky;
    top: 2%;
    background: white;
    justify-content: center;
    z-index: 99;
  }

  .main {
    height: 20rem;
    position: relative;
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    /* padding-bottom: 7rem; */
  }
  .logo {
  }
  main {
    text-align: center;
    padding: 1em;
    margin: 0 auto;
  }

  img {
    height: 8rem;
    width: 8rem;
  }

  h1 {
    color: #ff3e00;
    text-transform: uppercase;
    font-size: 4rem;
    font-weight: 100;
    line-height: 1.1;
    margin: 2rem auto;
    max-width: 14rem;
  }

  p {
    max-width: 14rem;
    margin: 1rem auto;
    line-height: 1.35;
  }

  @media (min-width: 480px) {
    h1 {
      max-width: none;
    }

    p {
      max-width: none;
    }
  }
</style>
