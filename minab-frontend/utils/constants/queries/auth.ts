import gql from "graphql-tag";


export const registerMutation = gql`
  mutation SignupUser(
    $first_name: String!, 
    $last_name: String!, 
    $email: String!, 
    $phone_number: String!, 
    $password: String!, 
    $remember_me: Boolean!
  ) {
    signup(
      first_name: $first_name, 
      last_name: $last_name, 
      email: $email, 
      phone_number: $phone_number, 
      password: $password, 
      remember_me: $remember_me
    ) {
      token
      user_id
      message
    }
  }
`;

export const loginMutation = gql`
  mutation LoginUser(
    $email: String!, 
    $password: String!, 
    $remember_me: Boolean!
  ) {
    # Call 'login' using the argument 'email'
    login(
      email: $email, 
      password: $password, 
      remember_me: $remember_me
    ) {
      token
      user_id
      message
    }
  }
`;