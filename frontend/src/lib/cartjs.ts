import { models } from "@wails/go/models";

export const automationByVendor: Record<
  string,
  ((selections: models.ProductSelection[]) => string) | undefined
> = {
  mtgmate: getMtgmateCartJs,
  gamescube: getGamescubeJs,
  goodgames: getGoodgamesJs,
};

// mtgmate
const mtgmateJs = `// 1. Open https://www.mtgmate.com.au in your browser
// 2. Open your browser tools/console and copy/paste the below block of code in.
//    You may need to do 'allow pasting'.
// ! This script cannot check whether all items were actually added to your cart
// ! Please review your cart and double check your cart quantities are correct

async function addToCart(sku, qty) {
    const formData = new FormData();
    formData.append("card_sku[id]", sku);
    formData.append("card_sku[quantity]", qty);
    const token_name = document.querySelector('head > meta[name="csrf-param"]').content;
    const token = document.querySelector('head > meta[name="csrf-token"]').content;
    formData.append(token_name, token);
    try {
        const resp = await fetch("https://www.mtgmate.com.au/cart/add_item_json", { method: "POST", body: formData });
        const data = await resp.json();
        if (!data.success) {
            console.error("Something went wrong adding item to cart. Response:")
            console.error(data);
        }
    } catch (err) {
        console.error(err);
    }
}
`;

function getMtgmateCartJs(selections: models.ProductSelection[]) {
  return (
    mtgmateJs +
    selections
      .map(
        (sel) =>
          `addToCart("${sel.offering.storeSKU}", "${sel.quantity.toString()}")`,
      )
      .join("\n")
  );
}

// gamescube
const gamescubeJs = `// 1. Open https://www.gamescube.com in your browser
// 2. Open your browser tools/console and copy/paste the below block of code in.
//    You may need to do 'allow pasting'.
// ! Please review your cart and double check your cart quantities are correct

async function addToCart(items) {
    const token = document.querySelector('head > meta[name="csrf-token"]').content;
    const payload = { line_items: items };
    try {
        const resp = await fetch("https://www.thegamescube.com/api/v1/cart/line_items", {
            method: "POST",
            body: JSON.stringify(payload),
            headers: {
                "x-csrf-token": token,
                "Content-Type": "application/json",
            },
        });
        const data = await resp.json();
        if (data.errors && data.errors.length > 0) {
            console.error("Something went wrong adding item to cart.");
            try {
                console.error(data.errors[0].message);
            } catch {
                console.error(data.errors);
            }
        }
    } catch (err) {
        console.error(err);
    }
}
const items = [];
`;

function getGamescubeJs(selections: models.ProductSelection[]) {
  return (
    gamescubeJs +
    selections
      .map(
        (sel) =>
          `items.push({ variant_id: ${sel.offering.storeSKU}, qty: ${sel.quantity.toString()} })`,
      )
      .join("\n") +
    "\naddToCart(items).then(() => console.log('ok'))"
  );
}

// goodgames
const goodgamesJs = `// 1. Open https://tcg.goodgames.com.au in your browser
// 2. Open your browser tools/console and copy/paste the below block of code in.
//    You may need to do 'allow pasting'.
// ! Please review your cart and double check your cart quantities are correct

async function addToCart(id, quantity) {
    try {
        const resp = await fetch("https://tcg.goodgames.com.au/cart/add.js", {
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ form_type: "product", utf8: "✓", quantity, id: 39894558572723 }),
            method: "POST",
        });
    } catch (err) {
        console.error(err);
    }
}
`;

function getGoodgamesJs(selections: models.ProductSelection[]) {
  return (
    goodgamesJs +
    selections
      .map(
        (sel) =>
          `addToCart(${sel.offering.storeSKU}, "${sel.quantity.toString()}"))`,
      )
      .join("\n")
  );
}
