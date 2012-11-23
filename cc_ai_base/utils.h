#ifndef UTILS_H_
#define UTILS_H_

#include "ai.h"

Card GetHighestCardInSuit(const CardSet& card_set, int suit);
Card GetLowestCardInSuit(const CardSet& card_set, int suit);

Card GetHighestCard(const CardSet& card_set);
Card GetLowestCard(const CardSet& card_set);

int GetWinnerIndex(const vector<Card>& played_cards);
Card GetHighestLoser(const CardSet& card_set,
                     const vector<Card>& played_cards);
Card GetLowestWinner(const CardSet& card_set,
                     const vector<Card>& played_cards);

#endif
