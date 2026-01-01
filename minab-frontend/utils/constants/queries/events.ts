import gql from "graphql-tag";

export const searchEventQuery = gql`
  query SearchEvent($take: Int = 10, $text: String = "%%") {
    events(where: { title: { _ilike: $text } }, limit: $take, order_by: { created_at: desc }) {
      id title featured_image event_date venue_name price category { name }
    }
  }
`;

export const GetMyEventsQuery = gql`
  query GetMyEvents($user_id: Int!) {
    events(where: { user_id: { _eq: $user_id } }, order_by: { created_at: desc }) {
      id title featured_image event_date venue_name category { name }
    }
    events_aggregate(where: { user_id: { _eq: $user_id } }) {
      aggregate { count }
    }
  }
`;

export const GetEventByIdQuery = gql`
  query GetEventById ($id: Int!, $user_id: Int!) {
    events_by_pk(id: $id) {
      id title description featured_image event_date price venue_name address
      images { image_url }
      event_tags { tag_name }
      category { name }
      # This checks if a bookmark exists for THIS user on THIS event
      event_bookmarks(where: { user_id: { _eq: $user_id } }) {
        user_id
      }
    }
  }
`;

export const GetReservedEventsQuery = gql`
  query GetReservedEvents($user_id: Int!) {
    tickets(where: { user_id: { _eq: $user_id } }) {
      id
      ticket_number
      status
      event {
        id
        title
        featured_image
        event_date
        venue_name
        category {
          name
        }
      }
    }
  }
`;

export const CREATE_EVENT_MUTATION = gql`
  mutation CreateEventAction($title: String!, $description: String!, $date: String!, $price: Float!, $location_lat: Float!, $location_lng: Float!, $venue_name: String!, $address: String!, $category_id: Int!, $tags: [String!]!, $image_urls: [String!]!, $featured_image: String!) {
    createEvent(title: $title, description: $description, date: $date, price: $price, location_lat: $location_lat, location_lng: $location_lng, venue_name: $venue_name, address: $address, category_id: $category_id, tags: $tags, image_urls: $image_urls, featured_image: $featured_image) { id message }
  }
`;

export const INSERT_TICKET_MUTATION = gql`
  mutation InsertTicket($event_id: Int!, $user_id: Int!, $ticket_number: String!) {
    insert_tickets_one(object: {
      event_id: $event_id, 
      user_id: $user_id, 
      status: "booked", 
      ticket_number: $ticket_number
    }) {
      id
    }
  }
`;

export const DELETE_BOOKMARK_MUTATION = gql`
  mutation DeleteB($event_id: Int!, $user_id: Int!) {
    delete_event_bookmarks(where: { _and: [{ event_id: { _eq: $event_id } }, { user_id: { _eq: $user_id } }] }) {
      affected_rows
    }
  }
`;
export const INSERT_BOOKMARK_MUTATION = gql`
  mutation InsertB($event_id: Int!, $user_id: Int!) {
    insert_event_bookmarks_one(object: {event_id: $event_id, user_id: $user_id}) { event_id }
  }
`;

export const INITIALIZE_PAYMENT_MUTATION = gql`
  mutation Pay($event_id: Int!, $full_name: String!, $email: String!) {
    initializePayment(event_id: $event_id, full_name: $full_name, email: $email) { checkout_url status message }
  }
`;
export const GetBookmarkedEventsQuery = gql`
  query GetBookmarkedEvents($user_id: Int!) {
    event_bookmarks(where: { user_id: { _eq: $user_id } }) {
      event {
        id
        title
        featured_image
        event_date
        venue_name
        category {
          name
        }
      }
    }
  }
`;