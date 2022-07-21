<script>
  import { fade, fly } from "svelte/transition";
  export let coins;

  export let selectedCoin;
  export let amountUSD;
  export let cryptoAmount;
  export let data;
  /**
   * @param {string} s
   */
  function recalc(s) {
    switch (s) {
      case "BTC":
        cryptoAmount === ""
          ? (amountUSD = 0.0)
          : (amountUSD = parseFloat(cryptoAmount) * 2000);
        break;
      case "USDT":
        cryptoAmount === ""
          ? (amountUSD = 0.0)
          : (amountUSD = parseFloat(cryptoAmount));
        break;
      case "ETH":
        cryptoAmount === ""
          ? (amountUSD = 0.0)
          : (amountUSD = parseFloat(cryptoAmount) * 500);
        break;
      case "USDC":
        cryptoAmount === ""
          ? (amountUSD = 0.0)
          : (amountUSD = parseFloat(cryptoAmount));
        break;
      default:
        break;
    }
  }
</script>

<div
  class="select_body"
  out:fly={{ x: 200, duration: 500 }}
  in:fly={{ x: -200, duration: 500 }}
>
  <div
    style="display: flex;
  flex-direction: column;
  align-items: center;"
  >
    <h3>Select Crypto Currency</h3>
    <div class="select-wrapper">
      <select
        class="select_form"
        bind:value={selectedCoin}
        on:change={(e) => {
          data.coin = selectedCoin;
          recalc(selectedCoin.symbol);
        }}
      >
        {#each coins as coin}
          <option value={coin}>
            {coin.text}
          </option>
        {/each}
      </select>
    </div>
  </div>
  <div class="v_divider" />
  <div>
    <h3>Select Donation Amount</h3>
    <div>
      <div class="usd_amount">{amountUSD.toFixed(2)}</div>
      <input
        class="select_form select"
        type="text"
        placeholder="0.00"
        bind:value={cryptoAmount}
        on:input={(e) => {
          recalc(selectedCoin.symbol);
          data.amount = {
            coin: cryptoAmount,
            fiat: amountUSD,
          };
          console.log(data);
        }}
      />
    </div>
  </div>
</div>

<style>
  select {
    -webkit-appearance: none;
    appearance: none;
  }
  .select-wrapper {
    position: relative;
  }

  .select-wrapper::after {
    content: "▼";
    font-size: 1rem;
    top: 1.2rem;
    right: 1rem;
    position: absolute;
  }
  .select_body {
    position: absolute;
    left: 0;
    width: 100%;
  }
  .select_form {
    height: auto;
    width: auto;
    font-size: larger;
    padding: 1rem;
    outline: none;
    border: 0;
    box-shadow: rgba(0, 0, 0, 0.05) 0px 6px 24px 0px,
      rgba(0, 0, 0, 0.08) 0px 0px 0px 1px;
  }
  .v_divider {
    margin-top: 3rem;
    margin-bottom: 3rem;
  }
  .usd_amount {
    margin-bottom: 1rem;
    font-size: xxx-large;
  }
  h3 {
    color: #ff3e00;
    text-transform: uppercase;
    font-size: 1rem;
    font-weight: 100;
    line-height: 1.1;
    margin: auto;
    margin-bottom: 0.5rem;
    margin-top: 1rem;
    max-width: 14rem;
  }
</style>
