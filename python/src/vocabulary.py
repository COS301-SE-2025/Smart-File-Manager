# to create a vocabulary of keywords

class Vocabulary:
    def __init__(self):
        self.vocab = []

    def createVocab(self, keywords):
        for filename, keywords in keywords.items():
            # print(f"\n== FILE: {filename} ==")
            for kw, score in keywords:
                self.vocab.append(kw)
                
    def getVocab(self):
        return self.vocab