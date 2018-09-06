using System;
using System.Collections.Generic;
using System.Text;

namespace Jargon
{
    public interface Tokens: IEnumerator<Token>
    {
        Token? Next();
    }
}
