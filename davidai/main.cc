#include "base.h"
#include "lasttwo.h"
#include "ravage.h"

#include <iostream>
#include <string>
#include <vector>
using namespace std;

void RunTrickTakingPlayer(AbstractTrickTakingPlayer* player) {
  while (true) {
    if (player->IsDone()) break;
    string line;
    getline(cin, line);
    vector<Card> initial_cards = Card::ListFromString(line);
    Card play;
    if (initial_cards.empty()) {
      play = player->LeadTrick();
    } else if (player->hand().GetValues(initial_cards[0].suit()).size()) {
      play = player->FollowTrick(initial_cards);
    } else {
      play = player->DiscardTrick(initial_cards);
    }
    cout << play.ToString() << endl;
    getline(cin, line);
    vector<Card> remaining_cards = Card::ListFromString(line);

    vector<Card> full_trick = initial_cards;
    full_trick.push_back(play);
    full_trick.insert(full_trick.end(), remaining_cards.begin(),
		      remaining_cards.end());
    player->RecordTrick(full_trick, initial_cards.size());
  }
}

int main() {
  while (true) {
    string game, hand;
    if (!getline(cin, game)) return 0;
    getline(cin, hand);
    CardSet cs = CardSet(hand);
    AbstractTrickTakingPlayer* player = NULL;
    if (game == "RAVAGE") {
      player = new RavagePlayer(cs);
    } else if (game == "LASTTWO") {
      player = new LastTwoPlayer(cs);
    } else {
      assert(false);
    }
    RunTrickTakingPlayer(player);
    delete player;
  }
}

/*int main() {
  while (true) {
    string line;
    getline(cin, line);
    CardSet cs = CardSet(line);
    AbstractTrickTakingPlayer* player = NULL;
    player = new RavagePlayer(cs);
    RunTrickTakingPlayer(player);
    delete player;
    if (!getline(cin, line)) return 0;
    assert(line == "RESET");
  }
}
*/
