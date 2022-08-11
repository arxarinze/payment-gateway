<script defer lang="ts">
  export let plugin_id :any;
  export let data: any;
  import QrCode from "qrcode";
  import { loadImage } from "canvas";
  import { onMount } from "svelte";
  import { fly } from "svelte/transition";
  import axios from "axios";
  const apiUrl = import.meta.env.VITE_API_URL;
  onMount(async () => {
    console.log(plugin_id);
    let address = await axios.post(
      apiUrl+"/api/v1/public/address",
      {
        crypto_symbol: data.coin.symbol,
        network: data.coin.network,
        plugin_id: plugin_id
      }
    );
    console.log("love1", address.data);
    let canvas = document.getElementById("canvas");
    console.log(data);
    QrCode.toCanvas(
      canvas,
      address?.data?.address!,
      {
        errorCorrectionLevel: "H",
        margin: 1,
        color: {
          dark: "#000000",
          light: "#ffffff",
        },
      },
      (error) => {
        //console.log(error, 'dsadas')
      }
    );
    const ctx = (canvas as HTMLCanvasElement)?.getContext("2d");
    var img: any;
    loadImage(
      "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAgAAAAIACAMAAADDpiTIAAAAY1BMVEX/////PgD/PgD/PgD/PgD/PgD/PgD/PgD/PgD/PgD/PgD/PgD/PgD/PgD/PgD/PgD/PgD/ShD/ViD/YjD/bkD/e1D/h2D/k3D/n4D/q5D/t6D/w7D/z8D/29D/6OD/9PD////+4eKDAAAAEHRSTlMAECAwQFBgcICQoLDA0ODwVOCoyAAAE31JREFUeNrs2lfCpSAMQOGANJWS/W92+utU/ssNzvnWQMlBBQAAAAAAAAAAAAAAAAAA4P/hQgi5fHeWH2IIh+DRQixnrforvd4lBy/PgiOV2vUvtLtELw8An86q/2bUEp3sCz5dXSf1K3nBfly8un6QfiUn2IjPt36wlg/BFnxu+hL9PATGuVT1hVgDtoVLX65lL7DI5a5r3FFgjb+GrtOLFxgSqq52BYERqeo71CQwIHV9l54Ebxa7LsAS4O63twTgb13A6BKAO9WKFmQ1pKGGVC8rwVc15nSyDIraM5KsgdDUpHoILA1/6xWB/ZcfemBj7lbjiuB18lDrtDEJ7N9+c7Jg5/abV51sjPabN4JsjvZjFqT9uAZovwl9+xqg/fg4QPsxCNB+Ey4nmHE03VtjBezffjwMv0no+gBjfgXQfqwA/vnc2Igyj/bbWZK/gjJUWQG0HyuA9mMF0H6f2TsTLUdx3wuTVCqVyobMFgMG+/1f8r/81pnpOX26c8HmGn0vUAs6tj5ZljOBNxNU91MbVPfTCFD30whQ99OTIXU/jQB1P4BrwY+6H8J38QPK4SrJqJt/oAWhhHyWEpumtb1zPvwJ5wbbNSZvFVD3a+w4h5/hnW2rNSPgUCiJWr7bfgq/xjx0lazEvVASuF/VjeH3mPp61URQOVwkEqabwjvMfc1/LqTu147hfWZbydKUxyIW6n7mNQeQsWFNA9T9KuvDAswd450xdb9qCEvhrZEl+VD3Wx1jA85aIfA8qPutzMuHhZk7dUEa92vmsAJTI4txUvdbDzOElRir7W8Ceuzf+rAa/qUmsHH3q1xYFVepCQAv/BAkf7EWgbu63/LUU4jAaHSk7G/ycZcY2BAH32hrwG+xUfcDeGmP6Obcrw8xGYzmgfzuhzAZzQM35H5jiM5caz2Q3/0QfC0oz2IB1P1cSIPv+O8JqPtBdLoE8LsfRKdLAL/7QXS6BPC7H0QjGGd94mVT7oe7gC4BB2L3wyNAlwB+9wNqgroEHFndD2CKfVVM3W8Km2IQhLs+75ne/UC6iIeC6n5z2B51pL4AdT8zhi0yGwE4FADqfvxpwJe6H9rynZ6W2QTV/XC82W9jyFfJ7344I5IG6u4Pu1962n12iJ/K9O7HvwmctfQHu196enmbW8HJ4Z6R++FUQClAl3/Y/dLjdrYHnDNzP5xmV3vAleC6b2TcjvaAw0Pdb8ljwbMO+8JbvtMz72UP+CjV/ZatBpV65fPPdD4w4nZxHnBW91teBL71+/8X6wMrg7zJQ7//FtwPp8pdBM/qfj/FZt4eflb3W8kEv/X7J3A/N9r/Y3Q+vQk+9PtHdr/xVf/BO7rBJ04D+b8/k/vNL/Nj/LmwAD7fSsAHh/sB430al3APuPDXf1ncrzdrDp0YMj0OODw38cIPztysm4f6PK8HHB4c7oeL+ivVTcG99n+YIcTD1b8gI3iI5ZcFXrhGPQGvPeAR4DK8I/jJ5n7Ao094BOSnAccyg5Zv38bqRW2ymxXyyNv9/gYXAGxuGvCdu/v9SOUTXBTlTwB43W9RGZzzmhZ0LHfgfj+CxGdeHnjPzP0imECT0+WAL373ew8gRruMPPBYsrvfu/TRNeCbbwMgcL+3qaMfCN73sAGkd78IaaDLJgCOJbP7GYEY4PnR/G2BN96W76kWkFd0DyQqARG4H0yjAfBkfeLFVYJT7z4ALtTuh4M3BXEHwKHkfOJlMJI6AJosasFXYvfD0QA4UpZ+rBENgG3VAG0a98Np9h0AJ1kE43jcTwNg+QWg8ancD6flD4DkC0CXzv1wLH8ApF4AhoTuh1O93H5bgk7pt3/c/XBMN/h9BsBVYMy0BffDaQe/vwA4bvD74+4HxMDe2oKvPN/fW4mA6aYItwMzOgUw08bcD6fq/W5uhnyxfH/fCc7yy4BlD4AnyfcfjUSmGcMv8JL3yMYB3WbdD6caVqsDPXJJAYf0Ld9pkwHD3RV8EIwXgfthGOvXGBN1z2MeZMPvfuDlVkc+IuQhCMaTux+eC1juADgKhCNxP5zGLTwr9JxDEcDSux/+qLXhPgp4CEBN7374mJtJ3uSYwQ4w0bsfPuuglzfJYAew/O6Hv3LQyHs8+HeAit/98HdOvLzJnX8HcPzuh197HeVNLvRVoJan5XvFRaATwAK5JwLM+bgfcPnVCHU/ULnFDHBuZfPUE7gDCPu7QMZn4H74/JNOqCXga3sLwNQICa1HdoArw32QBAuAFR6qCdgBvshTgC4D98OxnbzLB3kKMPO7X1rIU4CW3/3ScifvBnT87peWC/dBQJW/++EQNAPgDsjvfokgvxAwE7ufpgB4Dljzu19ivrifB+vV/UA+uOuAs7ofxrOgrgPWiUc98fPNfSfMqvuBfHJLgEvb8s1Pyf0+hFH3A7lyvw/QLP+8p+4ATBZo1f0wSvLpoE5LPxjf5AHgA8hLYlI1f8BoFei/pDoJ7CQKdWsHN4UfcK63TaVVoKJI1AsSQf6qrnfh53hn24r6HAAnTQ7Yrf7xhzn8IvPQGYnOgXw63LBh+6/tFH6T6VVRFgFwTikkwMl6VHYObxE3Bk7sAeADQCVr0WGB2UkkHgV7AASAfsVJfiC+ryQGZ+IR0bgFmhVHuOGMDZED4lzinwQMET4/hGsI2sGJA6Be65IuTQiUhz0HwCyLYz3X/cRLQR8AdkMpYDOTDSgoD7sOgFYWxQx002kvxa4DoFptVhfJPlAe9h0Ai4/sZBtTdil2HQDzClO6qBaB8rDvAHBLjuiKgG83eQ6sAWBsiES/4SIgzidBAGDuhzNVmz0GxDnFLgTNCSp/KL4muBHOEgAhpful72A6agDUAlKNIT52mwqIc4x+HNzh7peCIccMEGgKTfWfrF1IxCA4p3wCwKNZIGAftBFwLTbIM3pTaI27XxqG7BQAuBoGJGID7n6JGHLcAm7x12KDu18ihgyXgEv8q2EWd79UDARLAEFbcMXifovfab4VOBlcDBhx90tGl1klsDimuBvWpXI/HF/nZoIpjNzXuPulwptN3QrGeaQYEDAZ3P1SMQnCJRcPNCFGBLRz2B69ADxz8UCZ1o8AM4ZN0mZwLxSfFNnDfTaA+xGnAfdcNKBeud2ycmGzuKxMsMQOBAF6A1gGbT3oO5PjIBkCjO8A90uJrzJKAy8pXw2cX0Z+oHVh64z0Q4LxLFB8WISxM/IHmn4OBLT5VAMPgu8BINNoX03TtLZ3gYTZsF8Ow1+OrMOesfSlAPzdsCnsETwPvOWSBHRhzwz4iRB7EmB82DMVuwfgz0f3ugRk4QHfcd8P1yWgpG4Lw01Ql4BTJscBugRkch5w1SUgbmvIIwsR1CXAm32LoC4BXS4ieJN3qXzYMVMuScAZv7G9T+pMkoBDKe9i5kCC623bNM3LjlN6Ezxk4wHSBgLmvpE/YLoxLILPvysEnxWQHtfKD5jXnLIx5MI/KYQmD5wb+Xs6H2CGXLrDv+V9XqydG6ZPtgeUuVwPwK8KJx333fpUHnCk6Q4naAwABv5XU6Jy8GdBUApgN4Gx+oXgndLUgi40A+MIWkOAUf94BJhcpsVcBGEKicCf/TJTChF8FFvjIAjGM7jf31P7FP3hxea4CkLtebv2uwDg8p8XRdck7uqYxcx8ZgbeBOJF436pxx7+gy/KmYEEzSFjFflQuwU8cGPc+SPAt9FnXlkuDwSWAIII6E3833zU4fH/ZaBxv+WuOjugEEBkggQRYAVgjn4gSLUFEESAq1P93vkEwFNoI8C/BKOL7oHHTZ8GANikx/4JkoAmk0rQURaiS+d+ABoAgAKkPRcYjCzAvPcAOMlyVFMi9wNwse3jRJABEnSIWCOkAXDhfz0q/TNvUy2LoAFwKGVhKpfA/QDmfQfARZbn5aO7H0AgCACGBQBYBHD3A6j2HQAXWYfOR3U/gDb7AAAWAADTx3Q/gH7XAXARgPj7gLeyOHP0SuCFawEAaFykrg+AOvAHAL4AUITAUMkKDLsOgKesTjNu+fOL2XVb8FliUPUe3/uNrEO/64sBD4lE5wKA62QtKigs2QPgJPGoXtOb3mcrWQ8sMNkD4CoAUWJgsrWsiQ0IA3kAHCQ6pht+WbvnoatkXbo0zcj858AQVWudDz/H9V0VQVASTYrjd0CcqrG9c+EHJjfYtpYodAGk4g6AD9kApvkDlcTkFUA8+aTAb6Em/UGFIw+AUqhJ37JiuSfGf8p+qV1YgIZ7TtxVdosNi2AoJ0XqDtDMYREc49NxugOYPvmd9IM6QDpan/zRmFKrQMkwY1iMWagt8EPdD6QXagv80tIPvgMw9wPd1f0wZuGWANkZzbSZuVTFFjip+4FU3M9HX2RPtHNYGifvctUUIL374bTCnQOq+6VKAeWDuwrQ2T27H/56vJTkVYAxTBW/+6F4I+9y4z4KNvh8hvTuh2OFPAV4YH2UPb/7QXhDngIU6HvhU83pfukXgJL7JMgAU7pI3A9XAIIU4LxAJ/1ouN0PAMmBztx1wBGf1ZXe/VCcABy464A+hLD9RcD6sCaIBj+4xwI04c/4TrZHPSV4noZsNshy0zRcRel+AJMgHLnfh5jwmd3pW75BauHfAU6LzlOaWzr3A7CSwQ5wXnimqmtkE3Q+rI0TcAegtkALjHCjcD+gBozvAAQB4PAJzundD6IRiDP5G1Een+SX3v0QXoJxIAgAYKam7ysC9wMYBONKPh6yAca5ErgfkACyjQdcc6bO2Gza/QAmIxjPYius2141dYbC/YDvzz4iGj8KxJOBJO4HMBsBKQ/FVoAsEFgGSN1voRaoa7EVPuI8rTK0vO4HrP8Eb8afYj2v51eNAWMD0/e/FvwB8AZ+6Ex69wNxRnBOuwyAf+BelSyNGUI0BlmAe0EfAE14G5t81BPAS5bgtOsAaBO7H4BvJP0CQB8ATfqWb3z7hzhpAPC534J7173YdwBQut+SF+COGgB07rdk6notNADY3C+4WgAinAIARH9kmc/9Fr38eimyCAAJb1OTud+y912exe4DoOFyP9/Kkpw0AF4c7oc3OAMTAQh6Aiegms7jfnMji1Iei81xj/3O8kzkfkaW5avIJgCGJFmg6QlKPylqgADf0Wfu9fzuR78B4FfDuviPrFYjgfsBG0AizvHf2u6Y3I9gAyBoCcLTwNrxuh++AfC3BaOlADb3w58IJbgZ8j6+InU/gO9iqzwTjN52/O4HjoPgLwQ0AcFyux+eAPB7oASIjt/9uBMA/L0Ih0fAxtyvEwCgCYDUA23A6NK7HzDrFu4C49eAOoC8qNwP4HEoNs0TGhMFMJrNuF9vZDWeh2Lb3AQ6EATwLYv7AZQfxSbBNaANOK7BZ72geCv83x/hhD8YAOA6+SvmNdO5HzAPMj2C7wEIfugq+Q/1awyByf34v3/xwPcAFO9G+3+MLsRlNKLfv/iWN5lDcnD30+9ffG7pLdbo7qffvzjIm1SBmKkW/f7o47FDYMVb0e+PJwFNIMVV+v2XqASIC4z4Ttam/CyYKCWpCRK4H0P9D+AmWZpgsofNHh8FF2d5l07dDzj/5RdBtixgamR9rgUfN9mFCFgR0fR/2T1AHLP7afqH7wFV4MC/JAL3Q8HJTd6mV/fD23+Z9wDj1f3+SXkqaDmU8jatut8/uB0KYq7yPqO6n0h5Lqg5Sa6bgJUY3I8FOU95n3bn7iflV0HPlwAMvO6nuz9eChAx847d7/lZZMFVAOrdup98H4o8+BCEbqfu9zgV2XAXhIG/5Zu/9IfxKRCTuh+/CSKY7USAq9X93uEsELVX9yPnmUMEjJW6H74E0LYH+YTup0uAdPt2P10CpNu3++kSIF3+7vdRZMxZUDq/a/fTJUBqn7P7HYvMOQlMPSdzP73yiXMXGDOp+/HyIThmyND9nqdiH1xlAV7qfvwd4hCNT//Ei7of3B0IYMYQgUndbwUesggvr+7HnwciVC5hy7e6H8BFFuLl1f0YOTwFIM6jT4O6X8R6IEAzqfsR8i3L0c3pkz942JduAgDG+mXHvBp1PxoTWD4Ehoql5VtNYPmNwNtK1P3icJeFaWAjmF9GonA9FMqxlKWpLLIMjK3gqPvhV8Ug6v69GJg6I5G4HIpfQF0QiIEp/B5jV8nyqPuBaQBA1Q1z+DWmvpXlUffDWwNQqtY6H37GPNpGYNT9AD5kZUzT2cG58Ce8c71tG8FR90M5SzSaf1BLMv63nfu6khAGoiDa8gJkOv9o93u9G4Pouinw5jCF+6j9sKkBX7Qfdr0+2u8Lruul0X4sYFb5AsLUk6D9WABf+WYBj//UExLtxwWhJfDINwug/VgA7ccCbt9+SJP2owZpPxZgt/0QOu3HnSHL7Qe3G28/VOPthzx1SbPITSAM2+0H13Q1I8kNoRpvP8Sp6+hRYPg0UOUeUGg/Lgzbbj9U4+2H2I23H6rx9kNotB83B0y3H9xmvP3gG8980gO2P/aAPPQMDi92MYEWxSQmwJUfJjCynABy02doWU4CcTd+7oevQx9n7l7OBunQxxjFyRnBl/6AH3+U80LYht5Ry06sYgO9eFkCwu3PBcdaRx8u70NvZOzJyXrgbzCCsWcv64JLtU39m7Zx8K/Bp3p0/YXRag6Ca/Gx1KMN/UprW41RcGkhxlTfKDFG/ukBAAAAAAAAAAAAAAAAAADAjBeeQtAc4aIZhAAAAABJRU5ErkJggg=="
    ).then((data) => {
      img = data;
      img.height = 50;
      img.width = 50;
      const center = (150 - 50) / 2;
      ctx.drawImage(
        img,
        (canvas as HTMLCanvasElement).width / 2 - img.width / 2,
        (canvas as HTMLCanvasElement).height / 2 - img.height / 2,
        50,
        50
      );
    });
  });
</script>

<div
  class="deposit_body"
  out:fly={{ x: 200, duration: 500 }}
  in:fly={{ x: -200, duration: 500 }}
>
  <div>
    <h3>
      Use the address below to donate <span>{data.amount.coin}</span>
      <span>{data.coin.symbol}</span>
    </h3>
  </div>
  <div class="container">
    <canvas id="canvas" class="canvas-size" width="500" height="500" />
  </div>
</div>

<style>
  .canvas-size {
    width: 200px !important;
    height: 200px !important;
  }
  .deposit_body {
    display: flex;
    flex-direction: column;
    align-items: center;
    position: absolute;
    left: 0;
    width: 100%;
  }
  .container {
    display: flex;
    justify-content: center;
    align-items: center;
    border: 1px solid gray;
    margin-bottom: 3rem;
    border-radius: 50%;
  }
  .select_form {
    height: auto;
    width: auto;
    font-size: larger;
    padding: 1rem;
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
