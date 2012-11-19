#include "base.h"
#include "lasttwo.h"
#include "ravage.h"

#include <sstream>
#include <iostream>
#include <string>
#include <vector>
using namespace std;

void RunTrickTakingPlayer(AbstractTrickTakingPlayer* player) {
  string doubling;
  getline(cin, doubling);  // This should be "DOUBLING"
  for (int i = 0; i < 4; i++) {
    getline(cin, doubling);
    if (doubling == "DOUBLE") {
      switch(player->seat()) {
        case 0:
        cout << endl;
        break;

        case 1:
        cout << "0 2 3" << endl;
        break;

        case 2:
        cout << "0 1 3" << endl;
        break;

        case 3:
        cout << "0 1 2" << endl;
        break;

      }
    }
  }

  string line;
  while (true) {
    getline(cin, line);
    if (line != "TRICK") {
      break;
    }

    // TODO: Shouldn't need this check anymore since the protocol tells you
    // when you're done.
    if (player->IsDone()) break;

    vector<Card> full_trick;
    int position;
    for (int i = 0; i < 4; i++) {
      getline(cin, line);
      if (line == "PLAY") {
        Card play;
        position = i;
        player->PrepareForTrick();
        if (full_trick.empty()) {
          play = player->LeadTrick();
        } else if (player->hand().GetValues(full_trick[0].suit()).size()) {
          play = player->FollowTrick(full_trick);
        } else {
          play = player->DiscardTrick(full_trick);
        }
        cout << play.ToString() << endl;
        full_trick.push_back(play);
      } else {
        int seat;
        string card;
        stringstream ss(line);
        ss >> seat >> card;
        full_trick.push_back(card);
      }
    }
    player->RecordTrick(full_trick, position);
  }
  if (line.substr(0, 3) == "END") {
    // Can get scores here if we want them    
  }
}

int main() {
  while (true) {
    string player_line;
    if (!getline(cin, player_line)) return 0;
    stringstream ss(player_line);
    ss >> player_line;
    int seat;
    ss >> seat;

    string hand;
    getline(cin, hand);

    string game;
    getline(cin, game);

    CardSet cs = CardSet(hand);
    AbstractTrickTakingPlayer* player = NULL;
    if (game == "RAVAGE") {
      player = new RavagePlayer(seat, cs);
    } else if (game == "LASTTWO") {
      player = new LastTwoPlayer(seat, cs);
    } else {
      assert(false);
    }
    RunTrickTakingPlayer(player);
    delete player;
  }
}
