import gql from "graphql-tag";

export const initializePaymentMutation = gql`
  mutation InitializePayment($event_id: Int!, $quantity: Int!) {
    initializePayment(event_id: $event_id, quantity: $quantity) {
      checkout_url
      reference
    }
  }
`;