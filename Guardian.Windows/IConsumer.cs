using System;
namespace Guardian
{
    interface IConsumer<T>
    {
        T Consume();
    }
}
