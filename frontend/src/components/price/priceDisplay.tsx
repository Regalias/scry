export interface PriceProps {
    price: number;
}

export const Price = ({price}: PriceProps) => {
    return <span>${(price / 100).toFixed(2)}</span>
}