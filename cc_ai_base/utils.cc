#include "utils.h"

#include <assert.h>

Card GetHighestCardInSuit(const CardSet& card_set, int suit) {
  const vector<int>& values = card_set.GetValues(suit);
  if (values.empty())
    return Card();
  else
    return Card(suit, values.back());
}

Card GetLowestCardInSuit(const CardSet& card_set, int suit) {
  const vector<int>& values = card_set.GetValues(suit);
  if (values.empty())
    return Card();
  else
    return Card(suit, values[0]);
}

Card GetHighestCard(const CardSet& card_set) {
  Card best_card;
  for (int suit = 0; suit < 4; ++suit) {
    Card card = GetHighestCardInSuit(card_set, suit);
    if (!card.IsValid()) continue;
    if (!best_card.IsValid() || card.value() > best_card.value())
      best_card = card;
  }
  return best_card;
}

Card GetLowestCard(const CardSet& card_set) {
  Card best_card;
  for (int suit = 0; suit < 4; ++suit) {
    Card card = GetLowestCardInSuit(card_set, suit);
    if (!card.IsValid()) continue;
    if (!best_card.IsValid() || card.value() < best_card.value())
      best_card = card;
  }
  return best_card;
}

int GetWinnerIndex(const vector<Card>& played_cards) {
  assert(!played_cards.empty());
  int index = 0;
  for (int i = 1; i < played_cards.size(); ++i)
    if (played_cards[i].suit() == played_cards[index].suit() &&
	played_cards[i].value() > played_cards[index].value())
      index = i;
  return index;
}

Card GetHighestLoser(const CardSet& card_set,
                     const vector<Card>& played_cards) {
  Card best_card, winner_card = played_cards[GetWinnerIndex(played_cards)];
  const vector<int>& values = card_set.GetValues(winner_card.suit());
  for (int i = 0; i < values.size(); ++i) {
    if (values[i] > winner_card.value()) continue;
    if (!best_card.IsValid() || values[i] < best_card.value())
      best_card = Card(winner_card.suit(), values[i]);
  }
  return best_card;
}

Card GetLowestWinner(const CardSet& card_set,
                     const vector<Card>& played_cards) {
  Card best_card, winner_card = played_cards[GetWinnerIndex(played_cards)];
  const vector<int>& values = card_set.GetValues(winner_card.suit());
  for (int i = 0; i < values.size(); ++i) {
    if (values[i] < winner_card.value()) continue;
    if (!best_card.IsValid() || values[i] > best_card.value())
      best_card = Card(winner_card.suit(), values[i]);
  }
  return best_card;
}
